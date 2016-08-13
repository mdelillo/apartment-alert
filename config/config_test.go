package config_test

import (
	"os"

	"github.com/mdelillo/apartment-alert/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("New", func() {
		Describe("SMTP", func() {
			var (
				savedSMTPHost      string
				savedSMTPPassword  string
				savedSMTPPort      string
				savedSMTPRecipient string
				savedSMTPSender    string
				savedSMTPUsername  string
				savedVCAPServices  string
			)

			BeforeEach(func() {
				savedSMTPHost = os.Getenv("SMTP_HOST")
				savedSMTPPassword = os.Getenv("SMTP_PASSWORD")
				savedSMTPPort = os.Getenv("SMTP_PORT")
				savedSMTPRecipient = os.Getenv("SMTP_RECIPIENT")
				savedSMTPSender = os.Getenv("SMTP_SENDER")
				savedSMTPUsername = os.Getenv("SMTP_USERNAME")
				savedVCAPServices = os.Getenv("VCAP_SERVICES")
			})

			AfterEach(func() {
				os.Setenv("SMTP_HOST", savedSMTPHost)
				os.Setenv("SMTP_PASSWORD", savedSMTPPassword)
				os.Setenv("SMTP_PORT", savedSMTPPort)
				os.Setenv("SMTP_RECIPIENT", savedSMTPRecipient)
				os.Setenv("SMTP_SENDER", savedSMTPSender)
				os.Setenv("SMTP_USERNAME", savedSMTPUsername)
				os.Setenv("VCAP_SERVICES", savedVCAPServices)
			})

			Context("when there are SMTP environment variables", func() {
				BeforeEach(func() {
					os.Setenv("SMTP_HOST", "env-host")
					os.Setenv("SMTP_PASSWORD", "env-password")
					os.Setenv("SMTP_PORT", "1234")
					os.Setenv("SMTP_RECIPIENT", "env-recipient")
					os.Setenv("SMTP_SENDER", "env-sender")
					os.Setenv("SMTP_USERNAME", "env-username")
				})

				It("uses the environment variables", func() {
					c, err := config.New()
					Expect(err).NotTo(HaveOccurred())
					Expect(c.SMTPHost).To(Equal("env-host"))
					Expect(c.SMTPPassword).To(Equal("env-password"))
					Expect(c.SMTPPort).To(Equal(1234))
					Expect(c.SMTPRecipient).To(Equal("env-recipient"))
					Expect(c.SMTPSender).To(Equal("env-sender"))
					Expect(c.SMTPUsername).To(Equal("env-username"))
				})

				Context("when there is a VCAP_SERVICES with sendgrid", func() {
					BeforeEach(func() {
						os.Setenv("VCAP_SERVICES", `{"sendgrid":[{"credentials": {"hostname": "sendgrid-host", "username": "sendgrid-username", "password": "sendgrid-password"}}]}`)
					})

					It("uses the sendgrid credentials", func() {
						c, err := config.New()
						Expect(err).NotTo(HaveOccurred())
						Expect(c.SMTPHost).To(Equal("sendgrid-host"))
						Expect(c.SMTPPassword).To(Equal("sendgrid-password"))
						Expect(c.SMTPPort).To(Equal(25))
						Expect(c.SMTPRecipient).To(Equal("env-recipient"))
						Expect(c.SMTPSender).To(Equal("env-sender"))
						Expect(c.SMTPUsername).To(Equal("sendgrid-username"))
					})
				})

				Context("when there is a VCAP_SERVICES without sendgrid", func() {
					BeforeEach(func() {
						savedVCAPServices = os.Getenv("VCAP_SERVICES")
						os.Setenv("VCAP_SERVICES", `{"some-other-service":[{"credentials": {}}]}`)
					})

					It("uses the environment variables", func() {
						c, err := config.New()
						Expect(err).NotTo(HaveOccurred())
						Expect(c.SMTPHost).To(Equal("env-host"))
						Expect(c.SMTPPassword).To(Equal("env-password"))
						Expect(c.SMTPPort).To(Equal(1234))
						Expect(c.SMTPRecipient).To(Equal("env-recipient"))
						Expect(c.SMTPSender).To(Equal("env-sender"))
						Expect(c.SMTPUsername).To(Equal("env-username"))
					})
				})
			})

			Context("when there is no sendgrid service and no environment variables", func() {
				BeforeEach(func() {
					os.Setenv("SMTP_HOST", "")
					os.Setenv("SMTP_PASSWORD", "")
					os.Setenv("SMTP_PORT", "")
					os.Setenv("SMTP_RECIPIENT", "")
					os.Setenv("SMTP_SENDER", "")
					os.Setenv("SMTP_USERNAME", "")
				})

				It("has empty settings", func() {
					c, err := config.New()
					Expect(err).NotTo(HaveOccurred())
					Expect(c.SMTPHost).To(BeEmpty())
					Expect(c.SMTPPassword).To(BeEmpty())
					Expect(c.SMTPPort).To(Equal(0))
					Expect(c.SMTPRecipient).To(BeEmpty())
					Expect(c.SMTPSender).To(BeEmpty())
					Expect(c.SMTPUsername).To(BeEmpty())
				})
			})

			Context("when there is an error parsing VCAP_SERVICES", func() {
				BeforeEach(func() {
					savedVCAPServices = os.Getenv("VCAP_SERVICES")
					os.Setenv("VCAP_SERVICES", `some-bad-vcap-services`)
				})

				It("returns an error", func() {
					_, err := config.New()
					Expect(err).To(MatchError("could not parse 'some-bad-vcap-services' as json"))
				})
			})

			Context("when SMTP_PORT cannot be converted to an integer", func() {
				BeforeEach(func() {
					os.Setenv("SMTP_PORT", "some-bad-port")
				})

				It("returns an error", func() {
					_, err := config.New()
					Expect(err).To(MatchError("could not convert 'some-bad-port' to an integer"))
				})
			})
		})
	})
})
