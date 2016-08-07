package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const STREETEASY_BASE_URL = "https://streeteasy.com"

type Listing struct {
	Url     string
	Title   string
	Price   int
	Address string
	NoFee   bool
}
type DataLayer struct {
	Listings map[string]SearchResultListing
}
type SearchResultListing struct {
	Price      int
	AddrStreet string
	AddrUnit   string
	Title      string
	NoFee      bool
}

func main() {
	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	listings, err := GetListings(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listings:\n%+v", listings)
}

func GetListings(reader io.Reader) ([]Listing, error) {
	listings := make([]Listing, 0)

	resultsBytes, err := ioutil.ReadAll(reader)
	regex := regexp.MustCompile(`dataLayer = \[(.*)\];`)
	matches := regex.FindStringSubmatch(string(resultsBytes))
	if len(matches) <= 1 {
		panic("Could not find dataLayer in html")
	}

	dataLayer := parseDataLayer(matches[1])

	results, err := goquery.NewDocumentFromReader(bytes.NewReader(resultsBytes))
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
			Url:     STREETEASY_BASE_URL + href,
			Title:   dataLayer.Listings[listingId].Title,
			Price:   dataLayer.Listings[listingId].Price,
			Address: dataLayer.Listings[listingId].AddrStreet + " " + dataLayer.Listings[listingId].AddrUnit,
			NoFee:   dataLayer.Listings[listingId].NoFee,
		}
		listings = append(listings, listing)
	})

	return listings, nil
}

type dataLayer struct {
	SearchResultsListings map[string]interface{} `json:"searchResultsListings"`
}

func parseDataLayer(data string) *DataLayer {
	d := &dataLayer{}
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		panic(err)
	}

	parsedDataLayer := &DataLayer{
		Listings: make(map[string]SearchResultListing),
	}
	for _, searchResultsListing := range d.SearchResultsListings {
		id := strconv.FormatFloat(searchResultsListing.(map[string]interface{})["id"].(float64), 'f', 0, 64)
		parsedDataLayer.Listings[id] = SearchResultListing{
			Price:      int(searchResultsListing.(map[string]interface{})["price"].(float64)),
			AddrStreet: searchResultsListing.(map[string]interface{})["addr_street"].(string),
			AddrUnit:   searchResultsListing.(map[string]interface{})["addr_unit"].(string),
			Title:      searchResultsListing.(map[string]interface{})["title"].(string),
			NoFee:      searchResultsListing.(map[string]interface{})["noFee"].(bool),
		}
	}

	return parsedDataLayer
}
