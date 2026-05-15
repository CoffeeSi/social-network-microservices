package main

import (
	"log"

	"github.com/CoffeeSi/social-network-microservices/api-gateway/internal/app"
	"github.com/CoffeeSi/social-network-microservices/api-gateway/internal/config"
)

func main() {
	cfg := config.Load()
	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
