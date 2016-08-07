package main_test

import (
	"os"

	. "github.com/mdelillo/apartment-alert"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Describe("GetListings", func() {
		It("returns listings", func() {
			body, err := os.Open("fixtures/streeteasy-results.html")
			Expect(err).NotTo(HaveOccurred())

			expectedListing := Listing{
				Url:     "https://streeteasy.com/building/240-west-16-street-new_york/5e?featured=1",
				Title:   "240 West 16th #5E",
				Price:   2275,
				Address: "240 west 16th #5E",
				NoFee:   true,
			}

			listings, err := GetListings(body)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(listings)).To(Equal(8))
			Expect(listings[0]).To(Equal(expectedListing))
		})
	})
})
