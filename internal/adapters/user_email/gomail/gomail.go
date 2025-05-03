package gomail

import (
	"github.com/flastors/jwt-auth-golang/config"
	"github.com/flastors/jwt-auth-golang/pkg/logging"

	"gopkg.in/gomail.v2"
)

type userEmail struct {
	config config.SMTPConfig
	logger logging.Logger
	dealer *gomail.Dialer
}

func NewUserEmail(config config.SMTPConfig) *userEmail {
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	return &userEmail{
		config: config,
		dealer: d,
	}
}

func (e *userEmail) Send(to, subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	err := e.dealer.DialAndSend(m)
	if err != nil {
		e.logger.Errorln(err)
	}
}
