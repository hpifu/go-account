package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type MailClient struct {
	Password string
	From     string
	Server   string
	Port     int
}

func (m *MailClient) Send(to, subject, body string) error {
	content := fmt.Sprintf(`From: %v
To: %v
Subject: %v
Content-Type: text/html; charset=UTF-8;

%v
`, m.From, to, subject, body)

	return smtp.SendMail(
		fmt.Sprintf("%v:%v", m.Server, m.Port),
		smtp.PlainAuth("", m.From, m.Password, m.Server),
		m.From,
		strings.Split(to, ";"),
		[]byte(content),
	)
}
