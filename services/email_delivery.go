package services

import (
	"bytes"
	"context"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"html/template"
)

type Delivery interface {
	Send(user *model.User, text string) error
}

type activationTemplateContext struct {
	User            *model.User
	ActivationToken string
}

type ActivationTokenDelivery struct {
	MailTemplate *template.Template
	Delivery     Delivery
}

func (d *ActivationTokenDelivery) SendActivationToken(_ context.Context, user *model.User, token string) error {
	var buf bytes.Buffer
	ctx := activationTemplateContext{
		User:            user,
		ActivationToken: token,
	}
	err := d.MailTemplate.Execute(&buf, ctx)
	if err != nil {
		return err
	}

	return d.Delivery.Send(user, buf.String())
}

type PassCodeDelivery struct {
	MailTemplate *template.Template
	Delivery     Delivery
}

func (d *PassCodeDelivery) SendCode(_ context.Context, user *model.User, code string) error {
	var buf bytes.Buffer
	ctx := passwordDeliveryContext{
		User: user,
		Code: code,
	}
	if err := d.MailTemplate.Execute(&buf, ctx); err != nil {
		return err
	}
	return d.Delivery.Send(user, buf.String())
}

type passwordDeliveryContext struct {
	User *model.User
	Code string
}
