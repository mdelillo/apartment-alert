package parser_test

import (
	"fmt"

	"github.com/mdelillo/apartment-alert/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("New", func() {
		Context("when the URL domain is streeteasy.com", func() {
			It("returns a streeteasy parser", func() {
				p, err := parser.New("http://streeteasy.com/some/search", nil)
				Expect(err).NotTo(HaveOccurred())

				switch t := p.(type) {
				case *parser.StreetEasy:
					Expect(p.(*parser.StreetEasy).URL).To(Equal("http://streeteasy.com/some/search"))
				default:
					Fail(fmt.Sprintf("Expected StreetEasyParser, got %+v", t))
				}
			})
		})

		Context("when the URL domain is craigslist.org", func() {
			It("returns a craigslist parser", func() {
				p, err := parser.New("http://newyork.craigslist.org/some/search", nil)
				Expect(err).NotTo(HaveOccurred())

				switch t := p.(type) {
				case *parser.Craigslist:
					Expect(p.(*parser.Craigslist).URL).To(Equal("http://newyork.craigslist.org/some/search"))
				default:
					Fail(fmt.Sprintf("Expected StreetEasyParser, got %+v", t))
				}
			})
		})

		Context("when the URL domain is unknown", func() {
			It("returns an error", func() {
				_, err := parser.New("http://some-bad-domain.io/some/search", nil)
				Expect(err).To(MatchError("no parser for some-bad-domain.io"))
			})
		})
	})
})
