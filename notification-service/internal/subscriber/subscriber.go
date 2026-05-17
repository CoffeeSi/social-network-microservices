package subscriber

import (
	"encoding/json"
	"log"

	"github.com/CoffeeSi/social-network-microservices/notification-service/internal/event"
	"github.com/CoffeeSi/social-network-microservices/notification-service/internal/mailer"
	"github.com/nats-io/nats.go"
)

type NotificationSubscriber struct {
	nats   *nats.Conn
	mailer *mailer.Mailer
}

func NewNotificationSubscriber(nats *nats.Conn, mailer *mailer.Mailer) *NotificationSubscriber {
	return &NotificationSubscriber{
		nats:   nats,
		mailer: mailer,
	}
}

func (s *NotificationSubscriber) Start(subject string) error {
	_, err := s.nats.QueueSubscribe(subject, "notification_group", func(msg *nats.Msg) {
		log.Printf("New event: %s", subject)

		var ev event.UserVerificationEvent
		if err := json.Unmarshal(msg.Data, &ev); err != nil {
			log.Printf("Unmarshal error %v", err)
			return
		}

		err := s.sendEmail(ev.Email, ev.Code)
		if err != nil {
			log.Printf("Can't send email %s: %v", ev.Email, err)
			return
		}

		log.Printf("Email is send %s", ev.Email)
	})

	if err != nil {
		return err
	}

	return nil
}
func (s *NotificationSubscriber) sendEmail(email, code string) error {
	err := s.mailer.SendVerificationEmail(email, code)
	if err != nil {
		return err
	}
	return nil
}
