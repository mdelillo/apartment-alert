package parser

import (
	"errors"
	"net/url"
	"strings"

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

	switch strings.Split(u.Host, ":")[0] {
	case "newyork.craigslist.org":
		return &Craigslist{Fetcher: f, URL: searchURL}, nil
	case "renthop.com":
		return &RentHop{Fetcher: f, URL: searchURL}, nil
	case "streeteasy.com":
		return &StreetEasy{Fetcher: f, URL: searchURL}, nil
	case "127.0.0.1":
		return &Localhost{Fetcher: f, URL: searchURL}, nil
	default:
		return nil, errors.New("no parser for " + u.Host)
	}
}
