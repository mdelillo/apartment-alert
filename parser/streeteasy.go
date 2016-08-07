package parser

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const BASE_URL = "https://streeteasy.com"

type StreetEasy struct {
	Fetcher Fetcher
}

//go:generate mockgen -package mocks -destination mocks/fetcher.go github.com/mdelillo/apartment-alert/parser Fetcher
type Fetcher interface {
	Fetch(url string) (html string, err error)
}

type searchResultListing struct {
	Price      int
	AddrStreet string
	AddrUnit   string
	Title      string
	NoFee      bool
}

type dataLayer struct {
	SearchResultsListings map[string]interface{} `json:"searchResultsListings"`
}

func (s *StreetEasy) GetListings(url string) ([]Listing, error) {
	listings := make([]Listing, 0)

	html, err := s.Fetcher.Fetch(url)
	if err != nil {
		panic(err)
	}

	regex := regexp.MustCompile(`dataLayer = \[(.*)\];`)
	matches := regex.FindStringSubmatch(html)
	if len(matches) <= 1 {
		panic("Could not find dataLayer in html")
	}

	searchResultListings := parseDataLayer(matches[1])

	results, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		panic(err)
	}

	results.Find("div.details-title>a:first-child").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			html, _ := s.Html()
			panic("No href in " + html)
		}

		listingId, exists := s.Attr("data-gtm-listing-id")

		listing := Listing{
			ID:      listingId,
			Url:     BASE_URL + href,
			Title:   searchResultListings[listingId].Title,
			Price:   searchResultListings[listingId].Price,
			Address: searchResultListings[listingId].AddrStreet + " " + searchResultListings[listingId].AddrUnit,
			NoFee:   searchResultListings[listingId].NoFee,
		}
		listings = append(listings, listing)
	})

	return listings, nil
}

func parseDataLayer(data string) map[string]searchResultListing {
	d := &dataLayer{}
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		panic(err)
	}

	searchResultListings := make(map[string]searchResultListing)
	for _, searchResultsListing := range d.SearchResultsListings {
		id := strconv.FormatFloat(searchResultsListing.(map[string]interface{})["id"].(float64), 'f', 0, 64)
		searchResultListings[id] = searchResultListing{
			Price:      int(searchResultsListing.(map[string]interface{})["price"].(float64)),
			AddrStreet: searchResultsListing.(map[string]interface{})["addr_street"].(string),
			AddrUnit:   searchResultsListing.(map[string]interface{})["addr_unit"].(string),
			Title:      searchResultsListing.(map[string]interface{})["title"].(string),
			NoFee:      searchResultsListing.(map[string]interface{})["noFee"].(bool),
		}
	}

	return searchResultListings
}
