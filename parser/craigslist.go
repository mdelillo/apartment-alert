package parser

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Craigslist struct {
	Fetcher Fetcher
	URL     string
}

func (c *Craigslist) GetListings() ([]Listing, error) {
	listings := make([]Listing, 0)

	u, err := url.Parse(c.URL)
	if err != nil {
		panic(err)
	}
	baseUrl := u.Scheme + "://" + u.Host

	listingsHtml, err := c.Fetcher.Fetch(c.URL)
	if err != nil {
		panic(err)
	}

	listingsDoc, err := goquery.NewDocumentFromReader(strings.NewReader(listingsHtml))
	if err != nil {
		panic(err)
	}

	listingsDoc.Find("div.rows>p.row").Each(func(i int, selection *goquery.Selection) {
		id, exists := selection.Attr("data-pid")
		if !exists {
			html, _ := selection.Html()
			panic("No data-pid in " + html)
		}

		href, exists := selection.Find("a.hdrlnk").Attr("href")
		if !exists {
			html, _ := selection.Html()
			panic("No href in " + html)
		}

		title := selection.Find("span#titletextonly").Text()

		price, err := strconv.Atoi(selection.Find("span.price").Text()[1:])
		if err != nil {
			panic(err)
		}

		listingHtml, err := c.Fetcher.Fetch(baseUrl + href)
		if err != nil {
			panic(err)
		}

		listingDoc, err := goquery.NewDocumentFromReader(strings.NewReader(listingHtml))
		if err != nil {
			panic(err)
		}

		address := listingDoc.Find("div.mapaddress").Text()

		listings = append(listings, Listing{
			ID:      id,
			Url:     baseUrl + href,
			Title:   title,
			Price:   price,
			Address: address,
			NoFee:   true,
		})
	})

	return listings, nil
}
