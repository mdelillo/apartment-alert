package parser

import (
	"errors"
	"net/url"

	"github.com/mdelillo/apartment-alert/fetcher"
)

type Parser interface {
	GetListings() ([]Listing, error)
}

//go:generate mockgen -package mocks -destination mocks/fetcher.go github.com/mdelillo/apartment-alert/parser Fetcher
type Fetcher interface {
	Fetch(url string) (html string, err error)
}

type Listing struct {
	ID      string
	Url     string
	Title   string
	Price   int
	Address string
	NoFee   bool
}

func New(searchURL string, f *fetcher.Fetcher) (Parser, error) {
	u, err := url.Parse(searchURL)
	if err != nil {
		return nil, err
	}

	switch u.Host {
	case "newyork.craigslist.org", "craigslist.org":
		return &Craigslist{Fetcher: f, URL: searchURL}, nil
	case "streeteasy.com":
		return &StreetEasy{Fetcher: f, URL: searchURL}, nil
	default:
		return nil, errors.New("no parser for " + u.Host)
	}
}
