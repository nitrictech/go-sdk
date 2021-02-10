package v1

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Nitric Client", func() {

	When("A new client is created", func() {
		It("Should establish a new connection", func() {
			client, _ := New()
			client.Close()

		})
	})
})
