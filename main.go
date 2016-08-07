package main

import (
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/mdelillo/apartment-alert/alerter"
	"github.com/mdelillo/apartment-alert/fetcher"
	"github.com/mdelillo/apartment-alert/parser"
)

func main() {
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}
	a := &alerter.Alerter{
		Emailer: &alerter.EmailAlerter{
			SMTPUsername:  os.Getenv("SMTP_USERNAME"),
			SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
			SMTPHost:      os.Getenv("SMTP_HOST"),
			SMTPPort:      smtpPort,
			SMTPSender:    os.Getenv("SMTP_SENDER"),
			SMTPRecipient: os.Getenv("SMTP_RECIPIENT"),
			SendMail:      smtp.SendMail,
		},
	}
	streeteasy := parser.StreetEasy{
		Fetcher: &fetcher.Fetcher{},
	}

	url := os.Args[1]

	seenListings, err := streeteasy.GetListings(url)
	if err != nil {
		panic(err)
	}

	for {
		listings, err := streeteasy.GetListings(url)
		if err != nil {
			panic(err)
		}
		if err := a.Alert(listings, seenListings); err != nil {
			panic(err)
		}
		seenListings = listings
		time.Sleep(time.Minute)
	}
}
