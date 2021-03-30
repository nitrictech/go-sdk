package queueclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueueclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Queueclient Suite")
}
