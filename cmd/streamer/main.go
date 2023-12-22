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

func main() {
	cfg, err := config.Load(cfgType, cfgName, cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	nc, err := nats.Connect(cfg.NatsConfig.URL)
	if err != nil {
		log.Fatal(err)
	}

	sc, err := stan.Connect(cfg.NatsConfig.ClusterID, cfg.NatsConfig.ClientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Message %d", i)
		err := sc.Publish(cfg.NatsConfig.Canal, []byte(message))
		if err != nil {
			log.Printf("Error publishing message: %v\n", err)
		} else {
			log.Printf("Published message: %s\n", message)
		}

		// Небольшая задержка между сообщениями
		time.Sleep(time.Second)
	}
}