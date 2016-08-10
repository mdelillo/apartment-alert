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
		remoteAddr net.Addr
		from       string
		to         []string
		data       []byte
	}

	var (
		command                  string
		httpServerAddress        string
		savedApartmentAlertSleep string
		savedSMTPHost            string
		savedSMTPPassword        string
		savedSMTPPort            string
		savedSMTPRecipient       string
		savedSMTPSender          string
		savedSMTPUsername        string
		sentEmails               []Email
	)

	BeforeEach(func() {
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
		firstRequestListings := []parser.Listing{
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
		secondRequestListings := []parser.Listing{
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
		}
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())
		httpServerAddress = listener.Addr().String()

		requestCount := 0
		go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if requestCount == 0 {
				response, err := json.Marshal(initialListings)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(w, string(response))
			} else if requestCount == 1 {
				response, err := json.Marshal(firstRequestListings)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(w, string(response))
			} else if requestCount == 2 {
				response, err := json.Marshal(secondRequestListings)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(w, string(response))
			} else {
				w.WriteHeader(400)
				fmt.Fprintf(w, "too many requests made")
			}
			requestCount++
		}))

		savedApartmentAlertSleep = os.Getenv("APARTMENT_ALERT_SLEEP")
		savedSMTPHost = os.Getenv("SMTP_HOST")
		savedSMTPPassword = os.Getenv("SMTP_PASSWORD")
		savedSMTPPort = os.Getenv("SMTP_PORT")
		savedSMTPRecipient = os.Getenv("SMTP_RECIPIENT")
		savedSMTPSender = os.Getenv("SMTP_SENDER")
		savedSMTPUsername = os.Getenv("SMTP_USERNAME")
		os.Setenv("APARTMENT_ALERT_SLEEP", "0")
		os.Setenv("SMTP_USERNAME", "some-username")
		os.Setenv("SMTP_PASSWORD", "some-password")
		os.Setenv("SMTP_RECIPIENT", "some-recipient")
		os.Setenv("SMTP_SENDER", "some-sender")

		command, err = gexec.Build(filepath.Join("github.com", "mdelillo", "apartment-alert"))
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.Setenv("APARTMENT_ALERT_SLEEP", savedApartmentAlertSleep)
		os.Setenv("SMTP_HOST", savedSMTPHost)
		os.Setenv("SMTP_PASSWORD", savedSMTPPassword)
		os.Setenv("SMTP_PORT", savedSMTPPort)
		os.Setenv("SMTP_RECIPIENT", savedSMTPRecipient)
		os.Setenv("SMTP_SENDER", savedSMTPSender)
		os.Setenv("SMTP_USERNAME", savedSMTPUsername)
	})

	It("sends an email when new listings are found", func() {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())
		smtpServerAddress := listener.Addr().String()

		sentEmails = []Email{}
		smtpdServer := smtpd.Server{
			ConnectionChecker: func(peer smtpd.Peer) error { fmt.Println("Connection"); return nil },
			HeloChecker:       func(peer smtpd.Peer, name string) error { fmt.Println("Helo"); return nil },
			SenderChecker:     func(peer smtpd.Peer, addr string) error { fmt.Println("Sender"); return nil },
			RecipientChecker:  func(peer smtpd.Peer, addr string) error { fmt.Println("Recipient"); return nil },
			Authenticator: func(peer smtpd.Peer, username string, password string) error {
				fmt.Println("Auth")
				fmt.Printf("%+v\n", peer)
				fmt.Println(username)
				fmt.Println(password)
				return nil
			},
			Handler: func(peer smtpd.Peer, env smtpd.Envelope) error {
				sentEmails = append(sentEmails, Email{
					remoteAddr: peer.Addr,
					from:       env.Sender,
					to:         env.Recipients,
					data:       env.Data,
				})
				return nil
			},
		}
		fmt.Println("Imma listen on " + smtpServerAddress)
		smtpdServer.ListenAndServe(smtpServerAddress)
		os.Setenv("SMTP_HOST", strings.Split(smtpServerAddress, ":")[0])
		os.Setenv("SMTP_PORT", strings.Split(smtpServerAddress, ":")[1])

		_, err = gexec.Start(exec.Command(command, "http://"+httpServerAddress+"/search"), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(len(sentEmails), "100s").Should(Equal(2))
		Expect(sentEmails[0]).To(Equal(Email{
			remoteAddr: nil,
			from:       "some-sender",
			to:         []string{"some-recipient"},
			data:       []byte(""),
		}))
		Expect(sentEmails[0]).To(Equal(Email{
			remoteAddr: nil,
			from:       "some-sender",
			to:         []string{"some-recipient"},
			data:       []byte(""),
		}))
	})
})
