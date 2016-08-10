package parser

import "encoding/json"

type Localhost struct {
	Fetcher Fetcher
	URL     string
}

func (l *Localhost) GetListings() ([]Listing, error) {
	var listings []Listing

	listingsJSON, err := l.Fetcher.Fetch(l.URL)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(listingsJSON), &listings); err != nil {
		panic(err)
	}

	return listings, nil
}
