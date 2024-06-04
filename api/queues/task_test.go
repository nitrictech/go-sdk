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
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
)

var _ = Describe("ReceivedMessage", func() {
	var (
		ctrl        *gomock.Controller
		mockQ       *mock_v1.MockQueuesClient
		queueName   string
		leaseID     string
		message     map[string]interface{}
		receivedMsg ReceivedMessage
		ctx         context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockQ = mock_v1.NewMockQueuesClient(ctrl)
		queueName = "test-queue"
		leaseID = "1"
		message = map[string]interface{}{
			"message": "hello",
		}
		receivedMsg = &receivedMessageImpl{
			queueName:   queueName,
			queueClient: mockQ,
			leaseId:     leaseID,
			message:     message,
		}
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Message", func() {
		It("should return the correct message", func() {
			Expect(receivedMsg.Message()).To(Equal(message))
		})
	})

	Describe("Queue", func() {
		It("should return the correct queue name", func() {
			Expect(receivedMsg.Queue()).To(Equal(queueName))
		})
	})

	Describe("Complete", func() {
		It("should complete the message successfully", func() {
			mockQ.EXPECT().Complete(ctx, &v1.QueueCompleteRequest{
				QueueName: queueName,
				LeaseId:   leaseID,
			}).Return(&v1.QueueCompleteResponse{}, nil)

			err := receivedMsg.Complete(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle errors when completing the message", func() {
			mockQ.EXPECT().Complete(ctx, gomock.Any()).Return(nil, fmt.Errorf("some error"))

			err := receivedMsg.Complete(ctx)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Helper functions", func() {
	Describe("messageToWire", func() {
		It("should convert a map to a protobuf message", func() {
			message := map[string]interface{}{
				"message": "hello",
			}

			wireMsg, err := messageToWire(message)
			Expect(err).NotTo(HaveOccurred())
			Expect(wireMsg.GetStructPayload().AsMap()).To(Equal(message))
		})

		It("should handle errors when converting a map to a protobuf message", func() {
			message := map[string]interface{}{
				"message": make(chan int), // channels cannot be converted to protobuf
			}

			_, err := messageToWire(message)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("wireToMessage", func() {
		It("should convert a protobuf message to a map", func() {
			message := map[string]interface{}{
				"message": "hello",
			}
			wireMsg, err := messageToWire(message)
			Expect(err).NotTo(HaveOccurred())

			convertedMessage := wireToMessage(wireMsg)
			Expect(convertedMessage).To(Equal(message))
		})
	})
})