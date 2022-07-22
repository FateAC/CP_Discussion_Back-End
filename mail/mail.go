package mail

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpServer   = "smtp.gmail.com"
	smtpPort     = ":587"
	smtpUsername = "cpdiscussion.ta@gmail.com"
	smtpPassword = ""
)

func sendMail(to string, subject string, text string) error {
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	e := &email.Email{
		From:    "CPDiscussion <cpdiscussion.ta@gmail.com>",
		To:      []string{to},
		Subject: subject,
		Text:    []byte(text),
	}
	err := e.Send(smtpServer+smtpPort, auth)
	return err
}
