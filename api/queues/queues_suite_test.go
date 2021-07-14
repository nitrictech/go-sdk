package queues

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueuing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Queues Suite")
}
