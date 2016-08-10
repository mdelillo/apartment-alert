package parser_test

import (
	"encoding/json"

	"github.com/mdelillo/apartment-alert/parser"
	"github.com/mdelillo/apartment-alert/parser/mocks"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Localhost Parser", func() {
	var (
		localhostParser *parser.Localhost
		mockCtrl        *gomock.Controller
		mockFetcher     *mocks.MockFetcher
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFetcher = mocks.NewMockFetcher(mockCtrl)
		localhostParser = &parser.Localhost{
			Fetcher: mockFetcher,
			URL:     "http://127.0.0.1:1234/search",
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetListings", func() {
		It("returns listings", func() {
			expectedListings := []parser.Listing{
				{
					ID:      "id-1",
					Url:     "http://127.0.0.1:1234/listings/1",
					Title:   "title-1",
					Price:   1000,
					Address: "address-1",
					NoFee:   true,
				},
				{
					ID:      "id-2",
					Url:     "http://127.0.0.1:1234/listings/2",
					Title:   "title-2",
					Price:   2000,
					Address: "address-2",
					NoFee:   false,
				},
			}
			jsonListings, err := json.Marshal(expectedListings)
			Expect(err).NotTo(HaveOccurred())

			mockFetcher.EXPECT().Fetch("http://127.0.0.1:1234/search").Return(string(jsonListings), nil)

			listings, err := localhostParser.GetListings()
			Expect(err).NotTo(HaveOccurred())
			Expect(listings).To(Equal(expectedListings))
		})
	})
})
