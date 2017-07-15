package parser

import (
	"fmt"
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
		return nil, err
	}
	baseUrl := u.Scheme + "://" + u.Host

	listingsHtml, err := c.Fetcher.Fetch(c.URL)
	if err != nil {
		return nil, err
	}

	listingsDoc, err := goquery.NewDocumentFromReader(strings.NewReader(listingsHtml))
	if err != nil {
		return nil, err
	}

	selection := listingsDoc.Find("div.rows>p.row")
	for i := range selection.Nodes {
		node := selection.Eq(i)
		id, exists := node.Attr("data-pid")
		if !exists {
			html, err := node.Children().Html()
			if err != nil {
				panic(err)
			}
			return nil, fmt.Errorf("could not find data-pid attr in '%s'", html)
		}

		href, exists := node.Find("a.hdrlnk").Attr("href")
		if !exists {
			html, _ := node.Html()
			panic("No href in " + html)
		}

		title := node.Find("span#titletextonly").Text()

		price, err := strconv.Atoi(node.Find("span.price").Text()[1:])
		if err != nil {
			panic(err)
		}

		listingHtml, err := c.Fetcher.Fetch(baseUrl + href)
		if err != nil {
			panic(err)
		}

		listingDoc, err := goquery.NewDocumentFromReader(strings.NewReader(listingHtml))
		if err != nil {
			return nil, err
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
	}

	return listings, nil
}
