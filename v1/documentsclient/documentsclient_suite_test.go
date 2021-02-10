package documentsclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDocumentsclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Documentsclient Suite")
}
