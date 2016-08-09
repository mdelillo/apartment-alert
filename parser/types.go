package parser

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
