package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewMailer(host, port, username, password, from string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *Mailer) SendVerificationEmail(to, code string) error {
	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	subject := "Subject: Verification code"
	body := fmt.Sprintf("Your verification code is %s", code)
	msg := []byte(subject + "\r\n" + body)
	return smtp.SendMail(m.host+":"+m.port, auth, m.from, []string{to}, msg)
}
