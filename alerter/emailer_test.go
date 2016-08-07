package alerter_test

import (
	"errors"
	"net/smtp"

	"github.com/mdelillo/apartment-alert/alerter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmailAlerter", func() {
	var emailer *alerter.EmailAlerter

	BeforeEach(func() {
		emailer = &alerter.EmailAlerter{
			SMTPUsername:  "some-smtp-username",
			SMTPPassword:  "some-smtp-password",
			SMTPHost:      "some-smtp-host",
			SMTPPort:      12345,
			SMTPSender:    "some-smtp-sender",
			SMTPRecipient: "some-smtp-recipient",
		}
	})

	Describe("Send", func() {
		It("sends an email", func() {
			sendMail, recorder := mockSendMail(nil)
			emailer.SendMail = sendMail

			Expect(emailer.Send("some-msg")).To(Succeed())

			Expect(recorder.addr).To(Equal("some-smtp-host:12345"))
			Expect(recorder.auth).To(Equal(smtp.PlainAuth("", "some-smtp-username", "some-smtp-password", "some-smtp-host")))
			Expect(recorder.from).To(Equal("some-smtp-sender"))
			Expect(recorder.to).To(Equal([]string{"some-smtp-recipient"}))
			Expect(recorder.msg).To(Equal([]byte("To: some-smtp-recipient\nFrom: some-smtp-sender\nSubject: New apartment(s)\n\nsome-msg")))
		})

		Context("when sending an email fails", func() {
			It("returns an error", func() {
				sendMail, _ := mockSendMail(errors.New("some-error"))
				emailer.SendMail = sendMail

				Expect(emailer.Send("some-body")).To(MatchError("some-error"))
			})
		})
	})
})

func mockSendMail(err error) (func(string, smtp.Auth, string, []string, []byte) error, *emailRecorder) {
	recorder := new(emailRecorder)
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		*recorder = emailRecorder{addr, a, from, to, msg}
		return err
	}, recorder
}

type emailRecorder struct {
	addr string
	auth smtp.Auth
	from string
	to   []string
	msg  []byte
}
