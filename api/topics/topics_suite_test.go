package topics_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTopics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Topics Suite")
}
