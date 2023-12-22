package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TemaStatham/OrderService/pkg/cache"
	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const lifetimeElementInsideCache = 5 * time.Hour

// NatsConnConfig : a
type NatsConnConfig struct {
	URL       string
	ClientID  string
	ClusterID string
}

// StreamConnConfig :
type StreamConnConfig struct {
	Subject     string
	QueueGroup  string
	DurableName string
}

// Config :
type Config struct {
	NatsConnConfig
	StreamConnConfig
}

// Connect :
func Connect(c *cache.Cache, cfg Config) {
	nc, err := nats.Connect(cfg.NatsConnConfig.URL)
	if err != nil {
		log.Fatal(err)
		return
	}

	sc, err := stan.Connect(cfg.NatsConnConfig.ClusterID, cfg.NatsConnConfig.ClientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sc.Close()

	sub, err := sc.QueueSubscribe(
		cfg.StreamConnConfig.Subject,
		cfg.StreamConnConfig.QueueGroup,
		func(msg *stan.Msg) {
			handleRequest(msg, c)
		},
		stan.DurableName(cfg.StreamConnConfig.DurableName),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sub.Unsubscribe()

	waitForSignal(&sub)
}

func handleRequest(msg *stan.Msg, c *cache.Cache) {
	data := model.OrderClient{}

	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return
	}

	c.Set(data.OrderUID, data, lifetimeElementInsideCache)

	fmt.Printf("Received a message: %v\n", data)
}

func waitForSignal(sub *stan.Subscription) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Shutting down gracefully...")
}
