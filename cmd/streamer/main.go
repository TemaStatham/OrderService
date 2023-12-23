package main

import (
	"fmt"
	"log"
	"time"

	"github.com/TemaStatham/OrderService/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const (
	cfgName = "config"
	cfgType = "yaml"
	cfgPath = "./config"
)

const (
	publishTime = 5 * time.Second
	pubPref     = "pub"
)

func main() {
	cfg, err := config.Load(cfgType, cfgName, cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	opts := []nats.Option{nats.Name("NATS Streaming Example Publisher")}

	nc, err := nats.Connect(cfg.NatsConfig.URL, opts...)
	if err != nil {
		fmt.Print("asd\n\n\n\nasdasd\n")
		log.Fatal(err)
	}

	sc, err := stan.Connect(cfg.NatsConfig.ClusterID, cfg.NatsConfig.ClientID+pubPref, stan.NatsConn(nc))
	if err != nil {
		// fmt.Print("asd\n\n\n\n")
		log.Fatal(err)
	}
	defer sc.Close()

	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Message %d", i)
		err := sc.Publish(cfg.NatsConfig.Subject, []byte(message))
		if err != nil {
			log.Printf("Error publishing message: %v\n", err)
			continue
		}
		log.Printf("Published message: %s\n", message)

		time.Sleep(publishTime)
	}
}
