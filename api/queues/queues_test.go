package queues

import (
	"os"

	"github.com/golang/mock/gomock"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queues", func() {
	ctrl := gomock.NewController(GinkgoT())

	Context("New", func() {
		When("constructing a new queue client without the membrane", func() {
			os.Setenv("NITRIC_SERVICE_DIAL_TIMEOUT", "10")
			c, err := New()

			It("should return a nil client", func() {
				Expect(c).To(BeNil())
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		PWhen("constructing a new queue client without dial blocking", func() {
			// TODO: Do mock dial or non-blocking dial test here...
		})
	})

	Context("Queue", func() {
		When("creating a new Queue reference", func() {
			mockQ := mock_v1.NewMockQueueClient(ctrl)

			queues := &queuesImpl{
				c: mockQ,
			}

			q := queues.Queue("test-queue")

			qImpl, ok := q.(*queueImpl)

			It("Should have the provided queue name", func() {
				Expect(q.Name()).To(Equal("test-queue"))
			})

			It("Should be an instance of queueImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("Should share a client with the Queues client", func() {
				Expect(qImpl.c).To(Equal(mockQ))
			})
		})
	})
})
