package fetcher

import (
	"io/ioutil"
	"net/http"
)

type Fetcher struct{}

func (*Fetcher) Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(contents), nil
}
