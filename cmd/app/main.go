package main

import (
	"fmt"

	"github.com/TemaStatham/OrderService/config"
)

const (
	cfgName = "config"
	cfgType = "yaml"
	cfgPath = "./config"
)

func main() {
	cfg, err := config.Load(cfgType, cfgName, cfgPath)
	if err != nil {
		return
	}
	fmt.Print(*cfg)
}
