package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Fetcher struct{}

func (*Fetcher) Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %s", err)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response: %s: %s", resp.Status, contents)
	}

	return string(contents), nil
}
