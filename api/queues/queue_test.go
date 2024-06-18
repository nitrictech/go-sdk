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
	"errors"
	"strconv"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
)

var _ = Describe("Queue interface", func() {
	var (
		ctrl      *gomock.Controller
		mockQ     *mock_v1.MockQueuesClient
		queues    *queuesImpl
		queueName string
		q         Queue
		ctx       context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockQ = mock_v1.NewMockQueuesClient(ctrl)
		queues = &queuesImpl{
			queueClient: mockQ,
		}
		queueName = "test-queue"
		q = queues.Queue(queueName)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Having a valid queue", func() {
		Describe("Name", func() {
			It("should return the correct queue name", func() {
				Expect(q.Name()).To(Equal(queueName))
			})
		})

		Describe("Enqueue", func() {
			var messages []map[string]interface{}

			BeforeEach(func() {
				messages = []map[string]interface{}{
					{"message": "hello"},
					{"message": "world"},
				}
			})

			When("the operation is successful", func() {
				BeforeEach(func() {
					mockQ.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(
						&v1.QueueEnqueueResponse{
							FailedMessages: nil,
						},
						nil,
					).Times(1)
				})

				It("should successfully enqueue messages", func() {
					failedMessages, err := q.Enqueue(ctx, messages)
					Expect(err).ToNot(HaveOccurred())
					Expect(failedMessages).To(BeEmpty())
				})
			})

			When("a message send fails", func() {
				var failedMsg map[string]interface{}
				var failureReason string

				BeforeEach(func() {
					failedMsg = messages[0]
					wiredFailedMsg, err := messageToWire(failedMsg)
					Expect(err).ToNot(HaveOccurred())
					failureReason = "failed to send task"

					mockQ.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(
						&v1.QueueEnqueueResponse{
							FailedMessages: []*v1.FailedEnqueueMessage{
								{
									Message: wiredFailedMsg,
									Details: failureReason,
								},
							},
						},
						nil,
					).Times(1)
				})

				It("should recieve a message from []*FailedMessage", func() {
					failedMessages, err := q.Enqueue(ctx, messages)
					Expect(err).ToNot(HaveOccurred())
					Expect(failedMessages).To(HaveLen(1))
					Expect(failedMessages[0].Message).To(Equal(failedMsg))
					Expect(failedMessages[0].Reason).To(Equal(failureReason))
				})
			})

			When("the operation fails", func() {
				var errorMsg string

				BeforeEach(func() {
					errorMsg = "internal errror"
					mockQ.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(
						nil,
						errors.New(errorMsg),
					).AnyTimes()
				})

				It("should return an error", func() {
					_, err := q.Enqueue(ctx, messages)
					Expect(err).To(HaveOccurred())
					Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
				})
			})
		})

		Describe("Dequeue", func() {
			var depth int

			When("the depth is less than 1", func() {
				BeforeEach(func() {
					depth = 0
				})

				It("should return an error", func() {
					messages, err := q.Dequeue(ctx, depth)
					Expect(err).To(HaveOccurred())
					Expect(messages).To(BeNil())
				})
			})

			When("the operation is successful", func() {
				var messagesQueue []map[string]interface{}

				BeforeEach(func() {
					messagesQueue = []map[string]interface{}{
						{"message": "hello"},
						{"message": "world"},
					}
					depth = len(messagesQueue)

					dequeuedMessages := make([]*v1.DequeuedMessage, 0, depth)
					for i := 0; i < depth; i++ {
						msg, err := messageToWire(messagesQueue[i])
						Expect(err).ToNot(HaveOccurred())

						dequeuedMessages = append(dequeuedMessages, &v1.DequeuedMessage{
							LeaseId: strconv.Itoa(i),
							Message: msg,
						})
					}

					mockQ.EXPECT().Dequeue(ctx, &v1.QueueDequeueRequest{
						QueueName: queueName,
						Depth:     int32(depth),
					}).Return(&v1.QueueDequeueResponse{
						Messages: dequeuedMessages,
					}, nil).AnyTimes()
				})

				It("should receive tasks equal to depth", func() {
					messages, err := q.Dequeue(ctx, depth)
					Expect(err).ToNot(HaveOccurred())
					Expect(messages).To(HaveLen(depth))
				})

				It("should have message of type receivedMessageImpl", func() {
					messages, err := q.Dequeue(ctx, depth)
					Expect(err).ToNot(HaveOccurred())

					_, ok := messages[0].(*receivedMessageImpl)
					Expect(ok).To(BeTrue())
				})

				It("should contain the returned messages", func() {
					dequeuedMessages, err := q.Dequeue(ctx, depth)
					Expect(err).ToNot(HaveOccurred())

					for i := 0; i < depth; i++ {
						msg := dequeuedMessages[i]
						Expect(msg.Message()).To(Equal(messagesQueue[i]))
					}
				})
			})

			When("the operation fails", func() {
				var errorMsg string

				BeforeEach(func() {
					depth = 1
					errorMsg = "internal error"
					mockQ.EXPECT().Dequeue(gomock.Any(), gomock.Any()).Return(
						nil,
						errors.New(errorMsg),
					).Times(1)
				})

				It("should return an error", func() {
					_, err := q.Dequeue(ctx, depth)
					Expect(err).To(HaveOccurred())
					Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
				})
			})
		})
	})
})