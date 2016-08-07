package parser_test

import (
	"io/ioutil"

	"github.com/mdelillo/apartment-alert/parser"
	"github.com/mdelillo/apartment-alert/parser/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreetEasy Parser", func() {
	var (
		streetEasyParser *parser.StreetEasy
		mockCtrl         *gomock.Controller
		mockFetcher      *mocks.MockFetcher
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFetcher = mocks.NewMockFetcher(mockCtrl)
		streetEasyParser = &parser.StreetEasy{
			Fetcher: mockFetcher,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetListings", func() {
		It("returns listings", func() {
			expectedListings := []parser.Listing{
				{
					ID:      "1842962",
					Url:     "https://streeteasy.com/building/240-west-16-street-new_york/5e?featured=1",
					Title:   "240 West 16th #5E",
					Price:   2275,
					Address: "240 west 16th #5E",
					NoFee:   true,
				},
				{
					ID:      "1854006",
					Url:     "https://streeteasy.com/building/329-west-14-street-new_york/e1",
					Title:   "329 West 14th Street #E1",
					Price:   2250,
					Address: "329 West 14th Street #E1",
					NoFee:   true,
				},
				{
					ID:      "1851160",
					Url:     "https://streeteasy.com/building/102-greenwich-avenue-new_york/10",
					Title:   "102 Greenwich Avenue #10",
					Price:   2150,
					Address: "102 Greenwich Avenue  #10",
					NoFee:   true,
				},
				{
					ID:      "1844584",
					Url:     "https://streeteasy.com/building/172-5-avenue-manhattan/6a",
					Title:   "172 Fifth Avenue #6A",
					Price:   2350,
					Address: "172 Fifth Avenue #6A",
					NoFee:   true,
				},
				{
					ID:      "1838492",
					Url:     "https://streeteasy.com/building/240-west-16-street-new_york/d2",
					Title:   "240 W 16th #D2",
					Price:   2300,
					Address: "240 W 16TH #D2",
					NoFee:   true,
				},
				{
					ID:      "1832988",
					Url:     "https://streeteasy.com/building/240-west-16-street-new_york/b4",
					Title:   "240 West 16th Street #B4",
					Price:   2375,
					Address: "240 west 16th st #B4",
					NoFee:   true,
				},
				{
					ID:      "1828211",
					Url:     "https://streeteasy.com/building/307-west-29-street-new_york/1a",
					Title:   "307 W 29th Street #1A",
					Price:   2400,
					Address: "307 W 29th ST #1A",
					NoFee:   true,
				},
				{
					ID:      "1828212",
					Url:     "https://streeteasy.com/building/307-west-29-street-new_york/4a",
					Title:   "307 W 29th Street #4A",
					Price:   2308,
					Address: "307 W 29th ST #4A",
					NoFee:   true,
				},
			}

			contents, err := ioutil.ReadFile("fixtures/streeteasy-results.html")
			Expect(err).NotTo(HaveOccurred())

			mockFetcher.EXPECT().Fetch("some-url").Return(string(contents), nil)

			listings, err := streetEasyParser.GetListings("some-url")
			Expect(err).NotTo(HaveOccurred())
			Expect(listings).To(Equal(expectedListings))
		})
	})
})
