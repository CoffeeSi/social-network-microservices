package main

import (
	"log"

	"github.com/CoffeeSi/social-network-microservices/auth-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Error auth-service: %v", err)
	}
}
