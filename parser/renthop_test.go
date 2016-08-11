package parser_test

import (
	"io/ioutil"

	"github.com/mdelillo/apartment-alert/parser"
	"github.com/mdelillo/apartment-alert/parser/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RentHop Parser", func() {
	var (
		rentHopParser *parser.RentHop
		mockCtrl      *gomock.Controller
		mockFetcher   *mocks.MockFetcher
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFetcher = mocks.NewMockFetcher(mockCtrl)
		rentHopParser = &parser.RentHop{
			Fetcher: mockFetcher,
			URL:     "http://renthop.com/some-search",
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetListings", func() {
		It("returns listings", func() {
			expectedListings := []parser.Listing{
				{
					ID:      "7424119",
					Url:     "https://www.renthop.com/listings/w_34_st/7t/7424119",
					Title:   "Studio, 1BA at W 34 St.",
					Price:   2396,
					Address: "",
					NoFee:   true,
				},
				{
					ID:      "7427486",
					Url:     "https://www.renthop.com/listings/west_34th_street/8_t/7427486",
					Title:   "Studio, 1BA at West 34th Street",
					Price:   2330,
					Address: "",
					NoFee:   true,
				},
			}

			results, err := ioutil.ReadFile("fixtures/renthop-results.html")
			Expect(err).NotTo(HaveOccurred())

			mockFetcher.EXPECT().Fetch("http://renthop.com/some-search").Return(string(results), nil)

			listings, err := rentHopParser.GetListings()
			Expect(err).NotTo(HaveOccurred())
			Expect(listings).To(Equal(expectedListings))
		})
	})
})
