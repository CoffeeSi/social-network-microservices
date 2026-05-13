package event

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(url string) *Publisher {
	conn, err := nats.Connect(url)
	if err != nil {
		log.Printf("Failed to connect to NATS: %v", err)
		return nil
	}
	return &Publisher{conn: conn}
}

func (p *Publisher) Publish(subject string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.conn.Publish(subject, data)
}
