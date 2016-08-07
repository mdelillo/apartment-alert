package fetcher_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/mdelillo/apartment-alert/fetcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetcher", func() {
	var f *fetcher.Fetcher

	BeforeEach(func() {
		f = &fetcher.Fetcher{}
	})

	Describe("Fetch", func() {
		It("returns the contents of a URL", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "some-contents")
			}))
			defer server.Close()

			Expect(f.Fetch(server.URL)).To(Equal("some-contents"))
		})

		Context("when making the request fails", func() {
			It("returns an error", func() {
				_, err := f.Fetch("some-bad-url")
				Expect(err).To(MatchError("failed to fetch URL: Get some-bad-url: unsupported protocol scheme \"\""))
			})
		})

		Context("when the response code is not 200", func() {
			It("returns an error", func() {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
					fmt.Fprint(w, "some-error")
				}))
				defer server.Close()

				_, err := f.Fetch(server.URL)
				Expect(err).To(MatchError("invalid response: 404 Not Found: some-error"))
			})
		})
	})
})
