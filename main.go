package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mdelillo/apartment-alert/fetcher"
	"github.com/mdelillo/apartment-alert/parser"
)

func main() {
	p := parser.StreetEasy{
		Fetcher: &fetcher.Fetcher{},
	}
	listings, err := p.GetListings(os.Args[1])
	if err != nil {
		panic(err)
	}

	jsonListings, err := json.Marshal(listings)
	fmt.Println(string(jsonListings))
}
