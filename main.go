package main

import (
	"math/rand"
	"net/smtp"
	"os"
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

	for {
		listings, err := p.GetListings()
		if err != nil {
			panic(err)
		}
		if err := a.Alert(listings, seenListings); err != nil {
			panic(err)
		}
		seenListings = listings

		secondsToSleep := 600 + rand.Intn(60)
		time.Sleep(time.Duration(secondsToSleep) * time.Second)
	}
}
