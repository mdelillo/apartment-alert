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

		Context("when the URL domain is newyork.craigslist.org", func() {
			It("returns a craigslist parser", func() {
				p, err := parser.New("http://newyork.craigslist.org/some/search", nil)
				Expect(err).NotTo(HaveOccurred())

				switch t := p.(type) {
				case *parser.Craigslist:
					Expect(p.(*parser.Craigslist).URL).To(Equal("http://newyork.craigslist.org/some/search"))
				default:
					Fail(fmt.Sprintf("Expected CraigslistParser, got %+v", t))
				}
			})
		})

		Context("when the URL domain is renthop.com", func() {
			It("returns a renthop parser", func() {
				p, err := parser.New("http://renthop.com/some/search", nil)
				Expect(err).NotTo(HaveOccurred())

				switch t := p.(type) {
				case *parser.RentHop:
					Expect(p.(*parser.RentHop).URL).To(Equal("http://renthop.com/some/search"))
				default:
					Fail(fmt.Sprintf("Expected RentHopParser, got %+v", t))
				}
			})
		})

		Context("when the URL domain is 127.0.0.1", func() {
			It("returns a localhost parser", func() {
				p, err := parser.New("http://127.0.0.1:1234/search", nil)
				Expect(err).NotTo(HaveOccurred())

				switch t := p.(type) {
				case *parser.Localhost:
					Expect(p.(*parser.Localhost).URL).To(Equal("http://127.0.0.1:1234/search"))
				default:
					Fail(fmt.Sprintf("Expected LocalhostParser, got %+v", t))
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
