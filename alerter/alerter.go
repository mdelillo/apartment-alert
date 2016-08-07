package alerter

import (
	"fmt"

	"github.com/mdelillo/apartment-alert/parser"
)

type Alerter struct {
	Emailer Emailer
}

//go:generate mockgen -package mocks -destination mocks/emailer.go github.com/mdelillo/apartment-alert/alerter Emailer
type Emailer interface {
	Send(body string) error
}

func (a *Alerter) Alert(allListings []parser.Listing, seenListings []parser.Listing) error {
	emailBody := ""

	for _, listing := range allListings {
		isNewListing := true
		for _, seenListing := range seenListings {
			if listing.ID == seenListing.ID {
				isNewListing = false
				break
			}
		}

		if isNewListing {
			emailBody += fmt.Sprintf("%s $%d (%s) %s\n\n",
				listing.Title,
				listing.Price,
				formatNoFee(listing.NoFee),
				listing.Url,
			)
		}
	}

	if emailBody == "" {
		return nil
	}
	return a.Emailer.Send(emailBody)
}

func formatNoFee(noFee bool) string {
	if noFee {
		return "No Fee"
	} else {
		return "Fee"
	}
}
