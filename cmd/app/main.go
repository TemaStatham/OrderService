package main

import (
	"fmt"

	"github.com/TemaStatham/OrderService/config"
	"github.com/TemaStatham/OrderService/pkg/repository"
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

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     cfg.DB.DBHost,
		Port:     cfg.DB.DBPort,
		Username: cfg.DB.DBUser,
		Password: cfg.DB.DBPassword,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.DBSSLMode,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	repos := repository.NewRepository(db)
	fmt.Print(repos)
}
