package topicclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTopicclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Topicclient Suite")
}
