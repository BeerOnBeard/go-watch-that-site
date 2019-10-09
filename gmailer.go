package main

import (
	"fmt"
	"net/smtp"
)

const (
	// SMTPServer is the GMail SMTP server
	SMTPServer = "smtp.gmail.com"

	// SMTPServerPort is the SMTP server port
	SMTPServerPort = "587"
)

// GMailer sends HTML emails through GMail
type GMailer struct {
	User     string
	Password string
}

// Send an HTML email
func (mailer *GMailer) Send(to, subject, body string) error {
	message := fmt.Sprintf("From: %s\n", mailer.User)
	message += fmt.Sprintf("To: %s\n", to)
	message += fmt.Sprintf("Subject: %s\n", subject)
	message += "Content-Type: text/html; charset=\"UTF-8\";\n"
	message += "\n"
	message += body

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", SMTPServer, SMTPServerPort),
		smtp.PlainAuth("", mailer.User, mailer.Password, SMTPServer),
		mailer.User,
		[]string{to},
		[]byte(message))

	return err
}
