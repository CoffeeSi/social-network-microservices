package main

import (
	"log"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
