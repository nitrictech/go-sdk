package kvclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKVclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Documentsclient Suite")
}
