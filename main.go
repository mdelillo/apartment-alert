package main

import (
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/mdelillo/apartment-alert/alerter"
	"github.com/mdelillo/apartment-alert/config"
	"github.com/mdelillo/apartment-alert/fetcher"
	"github.com/mdelillo/apartment-alert/parser"
)

func main() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	a := &alerter.Alerter{
		Emailer: &alerter.EmailAlerter{
			Config:   config,
			SendMail: smtp.SendMail,
		},
	}
	p, err := parser.New(os.Args[1], &fetcher.Fetcher{})
	if err != nil {
		panic(err)
	}

	seenListings, err := p.GetListings()
	if err != nil {
		panic(err)
	}

	var secondsToSleep int
	if os.Getenv("APARTMENT_ALERT_SLEEP") != "" {
		secondsToSleep, err = strconv.Atoi(os.Getenv("APARTMENT_ALERT_SLEEP"))
		if err != nil {
			panic(err)
		}
	} else {
		secondsToSleep = 600
	}

	for {
		listings, err := p.GetListings()
		if err != nil {
			panic(err)
		}
		if err := a.Alert(listings, seenListings); err != nil {
			panic(err)
		}
		seenListings = listings

		time.Sleep(time.Duration(secondsToSleep) * time.Second)
	}
}
