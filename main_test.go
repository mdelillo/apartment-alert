package main_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bitbucket.org/chrj/smtpd"
	"github.com/mdelillo/apartment-alert/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	type Email struct {
		from string
		to   []string
		data []byte
	}

	var (
		command                  string
		savedApartmentAlertSleep string
		savedSMTPHost            string
		savedSMTPPort            string
		savedSMTPRecipient       string
		savedSMTPSender          string
		sentEmails               []Email
	)

	BeforeEach(func() {
		var err error
		command, err = gexec.Build(filepath.Join("github.com", "mdelillo", "apartment-alert"))
		Expect(err).NotTo(HaveOccurred())

		savedApartmentAlertSleep = os.Getenv("APARTMENT_ALERT_SLEEP")
		savedSMTPHost = os.Getenv("SMTP_HOST")
		savedSMTPPort = os.Getenv("SMTP_PORT")
		savedSMTPRecipient = os.Getenv("SMTP_RECIPIENT")
		savedSMTPSender = os.Getenv("SMTP_SENDER")
	})

	AfterEach(func() {
		os.Setenv("APARTMENT_ALERT_SLEEP", savedApartmentAlertSleep)
		os.Setenv("SMTP_HOST", savedSMTPHost)
		os.Setenv("SMTP_PORT", savedSMTPPort)
		os.Setenv("SMTP_RECIPIENT", savedSMTPRecipient)
		os.Setenv("SMTP_SENDER", savedSMTPSender)
	})

	It("sends an email when new listings are found", func() {
		initialListings := []parser.Listing{
			{
				ID:      "id-1",
				Url:     "url-1",
				Title:   "title-1",
				Price:   1000,
				Address: "address-1",
				NoFee:   true,
			},
			{
				ID:      "id-2",
				Url:     "url-2",
				Title:   "title-2",
				Price:   2000,
				Address: "address-2",
				NoFee:   false,
			},
		}
		firstUpdateListings := []parser.Listing{
			{
				ID:      "id-1",
				Url:     "url-1",
				Title:   "title-1",
				Price:   1000,
				Address: "address-1",
				NoFee:   true,
			},
			{
				ID:      "id-2",
				Url:     "url-2",
				Title:   "title-2",
				Price:   2000,
				Address: "address-2",
				NoFee:   false,
			},
			{
				ID:      "id-3",
				Url:     "url-3",
				Title:   "title-3",
				Price:   3000,
				Address: "address-3",
				NoFee:   true,
			},
		}
		secondUpdateListings := []parser.Listing{
			{
				ID:      "id-1",
				Url:     "url-1",
				Title:   "title-1",
				Price:   1000,
				Address: "address-1",
				NoFee:   true,
			},
			{
				ID:      "id-4",
				Url:     "url-4",
				Title:   "title-4",
				Price:   4000,
				Address: "address-4",
				NoFee:   false,
			},
			{
				ID:      "id-5",
				Url:     "url-5",
				Title:   "title-5",
				Price:   5000,
				Address: "address-5",
				NoFee:   false,
			},
		}

		httpListener, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())

		requestCount := 0
		go http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if requestCount == 0 {
				response, err := json.Marshal(initialListings)
				Expect(err).NotTo(HaveOccurred())
				fmt.Fprintf(w, string(response))
			} else if requestCount == 1 {
				response, err := json.Marshal(firstUpdateListings)
				Expect(err).NotTo(HaveOccurred())
				fmt.Fprintf(w, string(response))
			} else {
				response, err := json.Marshal(secondUpdateListings)
				Expect(err).NotTo(HaveOccurred())
				fmt.Fprintf(w, string(response))
			}
			requestCount++
		}))

		smtpListener, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())

		sentEmails = []Email{}
		smtpdServer := smtpd.Server{
			Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
				sentEmails = append(sentEmails, Email{
					from: env.Sender,
					to:   env.Recipients,
					data: env.Data,
				})
				return nil
			},
		}
		go smtpdServer.Serve(smtpListener)

		os.Setenv("SMTP_RECIPIENT", "some-recipient")
		os.Setenv("SMTP_SENDER", "some-sender")
		os.Setenv("SMTP_HOST", strings.Split(smtpListener.Addr().String(), ":")[0])
		os.Setenv("SMTP_PORT", strings.Split(smtpListener.Addr().String(), ":")[1])
		os.Setenv("APARTMENT_ALERT_SLEEP", "0")

		_, err = gexec.Start(exec.Command(command, "http://"+httpListener.Addr().String()+"/search"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() int { return len(sentEmails) }).Should(Equal(2))
		Expect(sentEmails[0]).To(Equal(Email{
			from: "some-sender",
			to:   []string{"some-recipient"},
			data: []byte("To: some-recipient\nFrom: some-sender\nSubject: New apartment(s)\n\ntitle-3 $3000 (No Fee) url-3\n\n"),
		}))
		Expect(sentEmails[1]).To(Equal(Email{
			from: "some-sender",
			to:   []string{"some-recipient"},
			data: []byte("To: some-recipient\nFrom: some-sender\nSubject: New apartment(s)\n\ntitle-4 $4000 (Fee) url-4\n\ntitle-5 $5000 (Fee) url-5\n\n"),
		}))
	})
})
