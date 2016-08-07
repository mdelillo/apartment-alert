package alerter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAlerter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alerter Suite")
}
