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

// import (
// 	"context"
// 	"fmt"

// 	"github.com/golang/mock/gomock"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"

// 	mock_v1 "github.com/nitrictech/go-sdk/mocks"
// 	v1 "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
// 	"github.com/nitrictech/protoutils"
// )

// var _ = Describe("Queue", func() {
// 	ctrl := gomock.NewController(GinkgoT())

// 	Context("Send", func() {
// 		When("the gRPC server returns an error", func() {
// 			mockQ := mock_v1.NewMockQueueServiceClient(ctrl)

// 			mockQ.EXPECT().SendBatch(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock error"))

// 			q := &queueImpl{
// 				name:        "test-queue",
// 				queueClient: mockQ,
// 			}

// 			_, err := q.Send(context.TODO(), []*Task{
// 				{
// 					ID:          "1234",
// 					PayloadType: "test-payload",
// 					Payload: map[string]interface{}{
// 						"test": "test",
// 					},
// 				},
// 			})

// 			It("should pass through the error", func() {
// 				Expect(err).To(HaveOccurred())
// 				Expect(err.Error()).To(Equal("Unknown: error from grpc library: \n mock error"))
// 			})
// 		})

// 		When("the task send succeeds", func() {
// 			mockQ := mock_v1.NewMockQueueServiceClient(ctrl)
// 			mockStruct, _ := protoutils.NewStruct(map[string]interface{}{
// 				"test": "test",
// 			})

// 			mockQ.EXPECT().SendBatch(gomock.Any(), gomock.Any()).Return(&v1.QueueSendBatchResponse{
// 				FailedTasks: []*v1.FailedTask{
// 					{
// 						Message: "Failed to send task",
// 						Task: &v1.NitricTask{
// 							Id:          "1234",
// 							PayloadType: "test-payload",
// 							Payload:     mockStruct,
// 						},
// 					},
// 				},
// 			}, nil)

// 			q := &queueImpl{
// 				name:        "test-queue",
// 				queueClient: mockQ,
// 			}

// 			fts, _ := q.Send(context.TODO(), []*Task{
// 				{
// 					ID:          "1234",
// 					PayloadType: "test-payload",
// 					Payload: map[string]interface{}{
// 						"test": "test",
// 					},
// 				},
// 			})

// 			It("should receive the failed tasks from the QueueSendBatchResponse", func() {
// 				Expect(fts).To(HaveLen(1))
// 				Expect(fts[0].Reason).To(Equal("Failed to send task"))
// 				Expect(fts[0].Task.ID).To(Equal("1234"))
// 				Expect(fts[0].Task.PayloadType).To(Equal("test-payload"))
// 				Expect(fts[0].Task.Payload).To(Equal(map[string]interface{}{
// 					"test": "test",
// 				}))
// 			})
// 		})

// 		Context("Receive", func() {
// 			When("Retrieving tasks with depth less than 1", func() {
// 				q := &queueImpl{
// 					name: "test-queue",
// 				}

// 				_, err := q.Receive(context.TODO(), 0)

// 				It("should return an error", func() {
// 					Expect(err).To(HaveOccurred())
// 					Expect(err.Error()).To(Equal("Invalid Argument: Queue.Receive: depth cannot be less than 1"))
// 				})
// 			})

// 			When("The grpc successfully returns", func() {
// 				mockStruct, _ := protoutils.NewStruct(map[string]interface{}{
// 					"test": "test",
// 				})
// 				mockQ := mock_v1.NewMockQueueServiceClient(ctrl)

// 				mockQ.EXPECT().Receive(gomock.Any(), gomock.Any()).Return(&v1.QueueReceiveResponse{
// 					Tasks: []*v1.NitricTask{
// 						{
// 							Id:          "1234",
// 							Payload:     mockStruct,
// 							PayloadType: "mock-payload",
// 							LeaseId:     "1234",
// 						},
// 					},
// 				}, nil)

// 				q := &queueImpl{
// 					name:        "test-queue",
// 					queueClient: mockQ,
// 				}

// 				t, _ := q.Receive(context.TODO(), 1)

// 				It("should receive a single task", func() {
// 					Expect(t).To(HaveLen(1))
// 				})

// 				rt, ok := t[0].(*receivedTaskImpl)

// 				It("the task should be of type recieveTaskImpl", func() {
// 					Expect(ok).To(BeTrue())
// 				})

// 				It("Should contain the returned task", func() {
// 					tsk := rt.Task()

// 					Expect(tsk.ID).To(Equal("1234"))
// 					Expect(tsk.PayloadType).To(Equal("mock-payload"))
// 					Expect(tsk.Payload).To(Equal(map[string]interface{}{
// 						"test": "test",
// 					}))
// 				})
// 			})
// 		})
// 	})
// })
