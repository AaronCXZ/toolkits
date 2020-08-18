package strategy

import "github.com/Muskchen/toolkits/email"

type Alerter interface {
	Send()
}

type Alert struct {
	Phones  string
	Mails   string
	Context string
	Sub     string
}

func NewAlerter(phone, mails, context, sub string) Alerter {
	return &Alert{
		Phones:  phone,
		Mails:   mails,
		Context: context,
		Sub:     sub,
	}
}

func (alert *Alert) Send() {
	go alert.sendMail()
}

func (alert *Alert) sendMail() {
	msg, err := email.NewMessage(alert.Mails, alert.Sub, alert.Context)
	if err != nil {
	}

	if err := email.Send(msg); err != nil {
	}
}
