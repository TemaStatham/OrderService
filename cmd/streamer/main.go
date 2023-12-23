package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	publishTime = 30 * time.Second
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
		order := GetRandomOrder()
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

// GetRandomOrder генерирует случайный заказ
func GetRandomOrder() model.OrderClient {
	rand.Seed(time.Now().UnixNano())

	return model.OrderClient{
		OrderUID:    getRandomString(10),
		TrackNumber: getRandomString(10),
		Entry:       getRandomString(10),
		Delivery: model.Delivery{
			Name:    getRandomString(10),
			Phone:   getRandomString(10),
			Zip:     getRandomString(10),
			City:    getRandomString(10),
			Address: getRandomString(10),
			Region:  getRandomString(10),
			Email:   getRandomString(10),
		},
		Payment: model.Payment{
			Transaction:  getRandomString(10),
			RequestID:    getRandomString(10),
			Currency:     getRandomString(5),
			Provider:     getRandomString(10),
			Amount:       rand.Intn(1000), // случайное значение int до 1000
			PaymentDt:    rand.Intn(1000), // случайное значение int до 1000
			Bank:         getRandomString(10),
			DeliveryCost: rand.Intn(1000), // случайное значение int до 1000
			GoodsTotal:   rand.Intn(1000), // случайное значение int до 1000
			CustomFee:    rand.Intn(100),  // случайное значение int до 100
		},
		Items: []model.Item{
			{
				ChrtID:      rand.Intn(1000),    // случайное значение int до 1000
				TrackNumber: getRandomString(10),
				Price:       rand.Intn(1000),    // случайное значение int до 1000
				Rid:         getRandomString(10),
				Name:        getRandomString(10),
				Sale:        rand.Intn(100),     // случайное значение int до 100
				Size:        getRandomString(5),
				TotalPrice:  rand.Intn(1000),    // случайное значение int до 1000
				NmID:        rand.Intn(1000),    // случайное значение int до 1000
				Brand:       getRandomString(10),
				Status:      rand.Intn(1000),    // случайное значение int до 1000
			},
		},
		Locale:            getRandomString(5),
		InternalSignature: getRandomString(10),
		CustomerID:        getRandomString(10),
		DeliveryService:   getRandomString(10),
		Shardkey:          getRandomString(5),
		SmID:              rand.Intn(1000),       // случайное значение int до 1000
		DateCreated:       getRandomTime(),        // случайная дата и время
		OofShard:          getRandomString(10),
	}
}

// getRandomString генерирует случайную строку заданной длины
func getRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// getRandomTime генерирует случайную дату и время
func getRandomTime() time.Time {
	year := rand.Intn(100) + 2000
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	hour := rand.Intn(24)
	minute := rand.Intn(60)
	second := rand.Intn(60)

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
}