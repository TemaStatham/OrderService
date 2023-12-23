package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/TemaStatham/OrderService/config"
	"github.com/TemaStatham/OrderService/pkg/model"
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
		log.Fatal(err)
	}

	sc, err := stan.Connect(cfg.NatsConfig.ClusterID, cfg.NatsConfig.ClientID+pubPref, stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for {
		order := GetConstOrder()
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Fatalf("Error marshaling order to JSON: %v", err)
		}
		err = sc.Publish(cfg.NatsConfig.Subject, []byte(orderJSON))
		if err != nil {
			log.Printf("Error publishing message: %v\n", err)
			continue
		}
		log.Printf("Published message")

		time.Sleep(publishTime)
	}
}

func GetConstOrder() model.OrderClient {
	return model.OrderClient{
		OrderUID:    "defaultOrderUID",
		TrackNumber: "defaultTrackNumber",
		Entry:       "defaultEntry",
		Delivery: model.Delivery{
			Name:    "defaultName",
			Phone:   "defaultPhone",
			Zip:     "defaultZip",
			City:    "defaultCity",
			Address: "defaultAddress",
			Region:  "defaultRegion",
			Email:   "defaultEmail",
		},
		Payment: model.Payment{
			Transaction:  "defaultTransaction",
			RequestID:    "defaultRequestID",
			Currency:     "defaultCurrency",
			Provider:     "defaultProvider",
			Amount:       0, // Установите значение по умолчанию для int
			PaymentDt:    0, // Установите значение по умолчанию для int
			Bank:         "defaultBank",
			DeliveryCost: 0, // Установите значение по умолчанию для int
			GoodsTotal:   0, // Установите значение по умолчанию для int
			CustomFee:    0, // Установите значение по умолчанию для int
		},
		Items: []model.Item{
			{
				ChrtID:      0, // Установите значение по умолчанию для int
				TrackNumber: "defaultTrackNumber",
				Price:       0, // Установите значение по умолчанию для int
				Rid:         "defaultRid",
				Name:        "defaultName",
				Sale:        0, // Установите значение по умолчанию для int
				Size:        "defaultSize",
				TotalPrice:  0, // Установите значение по умолчанию для int
				NmID:        0, // Установите значение по умолчанию для int
				Brand:       "defaultBrand",
				Status:      0, // Установите значение по умолчанию для int
			},
		},
		Locale:            "defaultLocale",
		InternalSignature: "defaultInternalSignature",
		CustomerID:        "defaultCustomerID",
		DeliveryService:   "defaultDeliveryService",
		Shardkey:          "defaultShardkey",
		SmID:              0,           // Установите значение по умолчанию для int
		DateCreated:       time.Time{}, // Установите значение по умолчанию для time.Time
		OofShard:          "defaultOofShard",
	}
}
