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

		secondsToSleep := 600 + rand.Intn(60)
		time.Sleep(time.Duration(secondsToSleep) * time.Second)
	}
}
