package services

import (
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"gopkg.in/gomail.v2"
	"os"
	"time"
)

type MailingConfig struct {
	FromAddress string
	FromName    string
	Subject     string
	ContentType string
}

type MailingService struct {
	Dialer *gomail.Dialer
	Cfg    MailingConfig
}

type TaskID int

func NewMailingService(dialer *gomail.Dialer, cfg MailingConfig) *MailingService {
	return &MailingService{
		Dialer: dialer,
		Cfg:    cfg,
	}
}

func (s *MailingService) Send(user *model.User, text string) error {
	msg := gomail.NewMessage()
	name := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	msg.SetAddressHeader("To", user.Email, name)
	msg.SetAddressHeader("From", s.Cfg.FromAddress, s.Cfg.FromName)
	msg.SetHeader("Subject", s.Cfg.Subject)
	msg.SetDateHeader("X-Date", time.Now())
	msg.SetBody(s.Cfg.ContentType, text)
	msg.WriteTo(os.Stdout)
	err := s.Dialer.DialAndSend(msg)
	return err
}
