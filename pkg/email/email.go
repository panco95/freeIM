package email

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

type Mail struct {
	Addr     string
	Identity string
	Username string
	Password string
	Host     string
}

func New(addr, identity, username, password, host string) *Mail {
	return &Mail{
		Addr:     addr,
		Identity: identity,
		Username: username,
		Password: password,
		Host:     host,
	}
}

func (m *Mail) SendEmail(subject, text string, toEmail []string) error {
	e := email.NewEmail()
	e.From = "IM通知 <1129443982@qq.com>"
	e.To = toEmail
	e.Subject = subject
	e.Text = []byte(text)
	err := e.Send(
		m.Addr,
		smtp.PlainAuth(
			m.Identity,
			m.Username,
			m.Password,
			m.Host,
		))
	return err
}
