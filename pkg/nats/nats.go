package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

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
func Connect(cfg Config) {
	nc, err := nats.Connect(cfg.NatsConnConfig.URL)
	if err != nil {
		log.Fatal(err)
	}

	sc, err := stan.Connect(cfg.NatsConnConfig.ClusterID, cfg.NatsConnConfig.ClientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()
	
	sub, err := sc.QueueSubscribe(cfg.StreamConnConfig.Subject, cfg.StreamConnConfig.QueueGroup, handleRequest, stan.DurableName(cfg.StreamConnConfig.DurableName))
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	waitForSignal(&sub)
}

func handleRequest(msg *stan.Msg) {
	data := model.OrderClient{}
	err := json.Unmarshal(msg.Data, &data)
	if err != nil {
		return
	}
	// if ok := s.addToCache(data); ok {
	// 	if err := s.addOrder(data); err != nil {
    //   log.Printf("Order adding error: %w\n", err)
    // }
	// 	log.Printf("Data are updated\n")
  	fmt.Printf("Received a message: %v\n", data)
}

func waitForSignal(sub *stan.Subscription) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Shutting down gracefully...")
}
