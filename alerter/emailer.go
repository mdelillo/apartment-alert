package alerter

import (
	"fmt"
	"net/smtp"

	"github.com/mdelillo/apartment-alert/config"
)

type EmailAlerter struct {
	Config   *config.Config
	SendMail func(string, smtp.Auth, string, []string, []byte) error
}

func (e *EmailAlerter) Send(body string) error {
	auth := smtp.PlainAuth(
		"",
		e.Config.SMTPUsername,
		e.Config.SMTPPassword,
		e.Config.SMTPHost,
	)
	return e.SendMail(
		fmt.Sprintf("%s:%d", e.Config.SMTPHost, e.Config.SMTPPort),
		auth,
		e.Config.SMTPSender,
		[]string{e.Config.SMTPRecipient},
		[]byte(fmt.Sprintf("To: %s\nFrom: %s\nSubject: New apartment(s)\n\n%s",
			e.Config.SMTPRecipient,
			e.Config.SMTPSender,
			body,
		)),
	)
}
