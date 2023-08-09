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
			mockQ := mock_v1.NewMockQueueServiceClient(ctrl)

			queues := &queuesImpl{
				queueClient: mockQ,
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
				Expect(qImpl.queueClient).To(Equal(mockQ))
			})
		})
	})
})
