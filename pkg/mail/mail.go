package mail

import (
	"fmt"
	"net/smtp"

	"github.com/Cyan903/c-share/pkg/log"
)

type MailClient struct {
	To       []string
	From     string
	Password string
	Host     string
	Port     int
}

func (m *MailClient) SendMail(msg []byte) error {
	auth := smtp.PlainAuth("", m.From, m.Password, m.Host)

	if err := smtp.SendMail(fmt.Sprintf("%s:%d", m.Host, m.Port), auth, m.From, m.To, msg); err != nil {
		log.Error.Println("Could not send email -", err)
		return err
	}

	return nil
}
