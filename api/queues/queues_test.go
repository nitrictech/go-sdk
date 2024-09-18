// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package queues

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("Queues API", func() {
	var (
		ctrl   *gomock.Controller
		mockQ  *mock_v1.MockQueuesClient
		queues *queuesImpl
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockQ = mock_v1.NewMockQueuesClient(ctrl)
		queues = &queuesImpl{
			queueClient: mockQ,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Queue method", func() {
		When("creating a new Queue reference", func() {
			var (
				q         Queue
				queueName string
				qImpl     *queueImpl
				ok        bool
			)

			BeforeEach(func() {
				queueName = "test-queue"
				q = queues.Queue(queueName)
				qImpl, ok = q.(*queueImpl)
			})

			It("should be an instance of queueImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provided queue name", func() {
				Expect(q.Name()).To(Equal(queueName))
			})

			It("should share the Queue's gRPC client", func() {
				Expect(qImpl.queueClient).To(Equal(mockQ))
			})
		})
	})

	Describe("New method", func() {
		When("constructing a new queue client without the membrane", func() {
			BeforeEach(func() {
				os.Setenv("NITRIC_SERVICE_DIAL_TIMEOUT", "10")
			})
			AfterEach(func() {
				os.Unsetenv("NITRIC_SERVICE_DIAL_TIMEOUT")
			})

			c, err := New()

			It("should return a nil client", func() {
				Expect(c).To(BeNil())
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		PWhen("constructing a new queue client without dial blocking", func() {
			// TODO:
		})
	})
})
