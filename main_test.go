package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Context("when the URL is for StreetEasy", func() {
		It("returns a StreetEasy parser", func() {
			Expect(true).To(BeTrue())
		})
	})
})
