package storageclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStorageclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storageclient Suite")
}
