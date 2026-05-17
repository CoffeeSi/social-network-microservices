package subscriber

import (
	"encoding/json"
	"log"

	"github.com/CoffeeSi/social-network-microservices/notification-service/internal/event"
	"github.com/CoffeeSi/social-network-microservices/notification-service/internal/mailer"
	"github.com/nats-io/nats.go"
)

type emailJob struct {
	email string
	code  string
}

type NotificationSubscriber struct {
	nats       *nats.Conn
	mailer     *mailer.Mailer
	jobQueue   chan emailJob
	numWorkers int
}

func NewNotificationSubscriber(nats *nats.Conn, mailer *mailer.Mailer, numWorkers int) *NotificationSubscriber {
	return &NotificationSubscriber{
		nats:       nats,
		mailer:     mailer,
		jobQueue:   make(chan emailJob, 100),
		numWorkers: numWorkers,
	}
}

func (s *NotificationSubscriber) Start(subject string) error {
	for i := range s.numWorkers {
		go s.worker(i)
	}

	_, err := s.nats.QueueSubscribe(subject, "notification_group", func(msg *nats.Msg) {
		log.Printf("New event: %s", subject)

		var ev event.UserVerificationEvent
		if err := json.Unmarshal(msg.Data, &ev); err != nil {
			log.Printf("Unmarshal error: %v", err)
			return
		}
		select {
		case s.jobQueue <- emailJob{email: ev.Email, code: ev.Code}:
		default:
			log.Printf("job queue is full, dropping email to %s", ev.Email)
		}
	})

	return err
}

func (s *NotificationSubscriber) worker(id int) {
	log.Printf("worker %d started", id)
	for job := range s.jobQueue {
		err := s.mailer.SendVerificationEmail(job.email, job.code)
		if err != nil {
			log.Printf("worker %d: can't send email to %s: %v", id, job.email, err)
			continue
		}
		log.Printf("worker %d: email sent to %s", id, job.email)
	}
}
