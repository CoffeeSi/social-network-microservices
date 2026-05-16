package main

import "github.com/CoffeeSi/social-network-microservices/notification-service/internal/app"

func main() {
	cfg := app.NewConfig()
	app.Run(cfg)
}
