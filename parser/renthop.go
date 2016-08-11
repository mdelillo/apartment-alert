package parser

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RentHop struct {
	Fetcher Fetcher
	URL     string
}

func (r *RentHop) GetListings() ([]Listing, error) {
	listings := make([]Listing, 0)

	listingsHtml, err := r.Fetcher.Fetch(r.URL)
	if err != nil {
		panic(err)
	}

	listingsDoc, err := goquery.NewDocumentFromReader(strings.NewReader(listingsHtml))
	if err != nil {
		panic(err)
	}

	listingsDoc.Find("div.search-listing").Each(func(i int, selection *goquery.Selection) {
		id, exists := selection.Attr("listing_id")
		if !exists {
			html, _ := selection.Html()
			panic("No listing_id in " + html)
		}

		url, exists := selection.Find("div.listing-search-text>a").Attr("href")
		if !exists {
			html, _ := selection.Html()
			panic("No url in " + html)
		}

		title := selection.Find("div.listing-search-text>a").Text()

		price, err := strconv.Atoi(strings.Replace(
			strings.Replace(
				strings.TrimSpace(
					selection.Find("div#listing-"+id+"-price").Text(),
				),
				",",
				"",
				1,
			),
			"$",
			"",
			1,
		))
		if err != nil {
			panic(err)
		}

		noFee := selection.Find(`div:contains("No Fee")`).Length() > 0

		listings = append(listings, Listing{
			ID:      id,
			Url:     url,
			Title:   title,
			Price:   price,
			Address: "",
			NoFee:   noFee,
		})
	})

	return listings, nil
}
