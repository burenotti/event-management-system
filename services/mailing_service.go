package services

import (
	"context"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"gopkg.in/gomail.v2"
	"time"
)

type MailingService struct {
	Dialer *gomail.Dialer
}

type TaskID int

func NewMailingService(dialer *gomail.Dialer) *MailingService {
	return &MailingService{
		Dialer: dialer,
	}
}

func (s *MailingService) SendMail(user *model.User, text string) error {
	msg := gomail.NewMessage()
	name := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	msg.SetAddressHeader("To", user.Email, name)
	msg.SetAddressHeader("From", "Shy0w1@yandex.ru", "Артём Буренин")
	msg.SetDateHeader("X-Date", time.Now())
	msg.SetBody("text/plain", text)

	err := s.Dialer.DialAndSend(msg)
	return err
}

func (s *MailingService) SendActivationToken(_ context.Context, user *model.User, token string) error {
	text := "Нажмите на ссылку, чтобы активировать ваш аккаунт http://localhost:8000/auth/activate/%s"
	return s.SendMail(user, fmt.Sprintf(text, token))
}

func (s *MailingService) SendCode(_ context.Context, user *model.User, code string) error {
	text := "Ваш код для входа: %s"
	return s.SendMail(user, fmt.Sprintf(text, code))
}
