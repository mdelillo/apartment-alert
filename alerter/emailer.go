package alerter

import (
	"fmt"
	"net/smtp"
)

type EmailAlerter struct {
	SMTPUsername  string
	SMTPPassword  string
	SMTPHost      string
	SMTPPort      int
	SMTPSender    string
	SMTPRecipient string
	SendMail      func(string, smtp.Auth, string, []string, []byte) error
}

func (e *EmailAlerter) Send(body string) error {
	fmt.Printf("email config: %+v\n", e)
	fmt.Printf("sending body: [%s]\n", body)
	auth := smtp.PlainAuth(
		"",
		e.SMTPUsername,
		e.SMTPPassword,
		e.SMTPHost,
	)
	return e.SendMail(
		fmt.Sprintf("%s:%d", e.SMTPHost, e.SMTPPort),
		auth,
		e.SMTPSender,
		[]string{e.SMTPRecipient},
		[]byte(fmt.Sprintf("To: %s\nFrom: %s\nSubject: New apartment(s)\n\n%s",
			e.SMTPRecipient,
			e.SMTPSender,
			body,
		)),
	)
}
