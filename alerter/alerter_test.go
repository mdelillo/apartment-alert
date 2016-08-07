package alerter_test

import (
	"github.com/golang/mock/gomock"
	"github.com/mdelillo/apartment-alert/alerter"
	"github.com/mdelillo/apartment-alert/alerter/mocks"
	"github.com/mdelillo/apartment-alert/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alerter", func() {
	var (
		a           *alerter.Alerter
		mockCtrl    *gomock.Controller
		mockEmailer *mocks.MockEmailer
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEmailer = mocks.NewMockEmailer(mockCtrl)
		a = &alerter.Alerter{
			Emailer: mockEmailer,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Alert", func() {
		It("sends alerts for new listings", func() {
			allListings := []parser.Listing{
				{
					ID:      "some-seen-id",
					Url:     "some-seen-url",
					Title:   "some-seen-title",
					Price:   1000,
					Address: "some-seen-address",
					NoFee:   true,
				},
				{
					ID:      "some-new-id",
					Url:     "some-new-url",
					Title:   "some-new-title",
					Price:   2000,
					Address: "some-new-address",
					NoFee:   true,
				},
				{
					ID:      "some-other-new-id",
					Url:     "some-other-new-url",
					Title:   "some-other-new-title",
					Price:   3000,
					Address: "some-other-new-address",
					NoFee:   false,
				},
			}
			seenListings := []parser.Listing{
				{
					ID:      "some-seen-id",
					Url:     "some-seen-url",
					Title:   "some-seen-title",
					Price:   1000,
					Address: "some-seen-address",
					NoFee:   true,
				},
			}
			expectedBody := `some-new-title
some-new-address
$2000 (No Fee)
some-new-url

some-other-new-title
some-other-new-address
$3000 (Fee)
some-other-new-url

`

			mockEmailer.EXPECT().Send(expectedBody)

			Expect(a.Alert(allListings, seenListings)).To(Succeed())
		})

		Context("when there are no new listings", func() {
			It("does not send an email", func() {
				allListings := []parser.Listing{
					{
						ID:      "some-seen-id",
						Url:     "some-seen-url",
						Title:   "some-seen-title",
						Price:   1000,
						Address: "some-seen-address",
						NoFee:   true,
					},
				}
				seenListings := []parser.Listing{
					{
						ID:      "some-seen-id",
						Url:     "some-seen-url",
						Title:   "some-seen-title",
						Price:   1000,
						Address: "some-seen-address",
						NoFee:   true,
					},
				}

				Expect(a.Alert(allListings, seenListings)).To(Succeed())
			})
		})
	})
})
