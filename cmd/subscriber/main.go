package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TemaStatham/OrderService/config"
	"github.com/TemaStatham/OrderService/pkg/cache"
	"github.com/TemaStatham/OrderService/pkg/handler"
	"github.com/TemaStatham/OrderService/pkg/repository"
	"github.com/TemaStatham/OrderService/pkg/server"
	"github.com/TemaStatham/OrderService/pkg/service"
)

const (
	cfgName = "config"
	cfgType = "yaml"
	cfgPath = "./config"

	cacheLifetime              = 0
	lifetimeElementInsideCache = 10 * time.Hour
)

func main() {
	cfg, err := config.Load(cfgType, cfgName, cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     cfg.DBConfig.DBHost,
		Port:     cfg.DBConfig.DBPort,
		Username: cfg.DBConfig.DBUser,
		Password: cfg.DBConfig.DBPassword,
		DBName:   cfg.DBConfig.DBName,
		SSLMode:  cfg.DBConfig.DBSSLMode,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	c := cache.New(cacheLifetime, lifetimeElementInsideCache)

	// nats.Connect(c, nats.Config{
	// 	NatsConnConfig: nats.NatsConnConfig{
	// 		URL:       cfg.NatsConfig.URL,
	// 		ClientID:  cfg.NatsConfig.ClientID,
	// 		ClusterID: cfg.NatsConfig.ClusterID,
	// 	},
	// 	StreamConnConfig: nats.StreamConnConfig{
	// 		Subject:     cfg.NatsConfig.Subject,
	// 		QueueGroup:  cfg.NatsConfig.QueueGroup,
	// 		DurableName: cfg.NatsConfig.DurableName,
	// 	},
	// })

	repos := repository.NewRepository(db, c)
	service := service.NewService(repos)
	hand := handler.NewHandler(service)
	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.ServConfig.Port, hand.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatal("error occured on db connection close: %s", err.Error())
	}
}
