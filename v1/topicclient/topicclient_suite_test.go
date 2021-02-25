package topicclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEventclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Eventclient Suite")
}
