package parser_test

import (
	"errors"
	"io/ioutil"

	"github.com/mdelillo/apartment-alert/parser"
	"github.com/mdelillo/apartment-alert/parser/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Craigslist Parser", func() {
	var (
		craigslistParser *parser.Craigslist
		mockCtrl         *gomock.Controller
		mockFetcher      *mocks.MockFetcher
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFetcher = mocks.NewMockFetcher(mockCtrl)
		craigslistParser = &parser.Craigslist{
			Fetcher: mockFetcher,
			URL:     "http://newyork.craigslist.org/some-search",
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetListings", func() {
		It("returns listings", func() {
			expectedListings := []parser.Listing{
				{
					ID:      "5714884495",
					Url:     "http://newyork.craigslist.org/mnh/abo/5714884495.html",
					Title:   "Bright... Short / Long Term ...NO FEE....SEE TODAY....",
					Price:   2249,
					Address: "240 west 16th stn at 8th Ave",
					NoFee:   true,
				},
				{
					ID:      "5691703979",
					Url:     "http://newyork.craigslist.org/mnh/abo/5691703979.html",
					Title:   "W/ GARDEN....see today..5PM.....NO FEE.... by owner",
					Price:   2199,
					Address: "430 east 13th at 1st ave",
					NoFee:   true,
				},
			}

			results, err := ioutil.ReadFile("fixtures/craigslist-results.html")
			Expect(err).NotTo(HaveOccurred())
			listing1, err := ioutil.ReadFile("fixtures/craigslist-listing-1.html")
			Expect(err).NotTo(HaveOccurred())
			listing2, err := ioutil.ReadFile("fixtures/craigslist-listing-2.html")
			Expect(err).NotTo(HaveOccurred())

			mockFetcher.EXPECT().Fetch("http://newyork.craigslist.org/some-search").Return(string(results), nil)
			mockFetcher.EXPECT().Fetch("http://newyork.craigslist.org/mnh/abo/5714884495.html").Return(string(listing1), nil)
			mockFetcher.EXPECT().Fetch("http://newyork.craigslist.org/mnh/abo/5691703979.html").Return(string(listing2), nil)

			listings, err := craigslistParser.GetListings()
			Expect(err).NotTo(HaveOccurred())
			Expect(listings).To(Equal(expectedListings))
		})

		Context("when there are no listings in the HTML", func() {
			It("returns no listings", func() {
				mockFetcher.EXPECT().Fetch("http://newyork.craigslist.org/some-search").Return(string("<html />"), nil)

				listings, err := craigslistParser.GetListings()
				Expect(err).NotTo(HaveOccurred())
				Expect(listings).To(BeEmpty())
			})
		})

		Context("when the URL is not valid", func() {
			It("returns an error", func() {
				craigslistParser.URL = "http://[::1]bad-port"
				_, err := craigslistParser.GetListings()
				Expect(err).To(MatchError(`parse http://[::1]bad-port: invalid port "bad-port" after host`))
			})
		})

		Context("when fetching the URL fails", func() {
			It("returns an error", func() {
				mockFetcher.EXPECT().Fetch("http://newyork.craigslist.org/some-search").Return("", errors.New("some-error"))
				_, err := craigslistParser.GetListings()
				Expect(err).To(MatchError("some-error"))
			})
		})

		FContext("when a listing does not have a data-pid", func() {
			It("returns an error", func() {
				html := `<html><body><div class="rows"><p class="row" data-pid="1" /><div></body></html>`
				mockFetcher.EXPECT().Fetch("http://newyork.craigslist.org/some-search").Return(html, nil)
				_, err := craigslistParser.GetListings()
				Expect(err).To(MatchError(`could not find data-pid attr in '<p class="row" />'`))
			})
		})
	})
})
