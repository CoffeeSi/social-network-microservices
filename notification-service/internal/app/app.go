package app

import (
	"log"

	"github.com/CoffeeSi/social-network-microservices/notification-service/internal/mailer"
	"github.com/CoffeeSi/social-network-microservices/notification-service/internal/subscriber"
	"github.com/nats-io/nats.go"
)

func Run(cfg *Config) {

	nc, err := nats.Connect(cfg.NatsURI)
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}
	defer nc.Close()

	m := mailer.NewMailer(
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.From,
	)

	sub := subscriber.NewNotificationSubscriber(nc, m)
	if err := sub.Start("user.verification"); err != nil {
		log.Fatalf("failed to start subscriber: %v", err)
	}

	log.Println("notification service started")
	select {}
}
