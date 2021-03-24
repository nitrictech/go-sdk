package queueclient

import (
	"fmt"

	"github.com/golang/mock/gomock"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ = Describe("Eventclient", func() {
	ctrl := gomock.NewController(GinkgoT())

	When("Push", func() {
		When("The queue exists", func() {
			When("1 event is pushed", func() {
				It("Should push the event", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Push with a single event")
					mockQueueClient.EXPECT().SendBatch(gomock.Any(), &v1.QueueSendBatchRequest{
						Queue: "test-queue",
						Tasks: []*v1.NitricTask{
							{
								Id:          "test-request-id",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
						},
					}).Return(&v1.QueueSendBatchResponse{}, nil)

					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{
						"hello": 123,
					}

					client := NewWithClient(mockQueueClient)
					res, err := client.SendBatch(&SendBatchOptions{
						Queue: "test-queue",
						Tasks: []*Task{
							&Task{
								ID:          requestId,
								PayloadType: payloadType,
								Payload:     payload,
							},
						},
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Not returning failed messages")
					Expect(res.FailedTasks).To(HaveLen(0))
				})
			})

			When("Multiple events are pushed", func() {
				It("Should push the events", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Push with multiple events")
					mockQueueClient.EXPECT().SendBatch(gomock.Any(), &v1.QueueSendBatchRequest{
						Queue: "test-queue",
						Tasks: []*v1.NitricTask{
							{
								Id:          "test-request-id",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
							{
								Id:          "test-request-id2",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
						},
					}).Return(&v1.QueueSendBatchResponse{}, nil)

					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{
						"hello": 123,
					}
					requestId2 := "test-request-id2"

					client := NewWithClient(mockQueueClient)
					res, err := client.SendBatch(&SendBatchOptions{
						Queue: "test-queue",
						Tasks: []*Task{
							{
								ID:          requestId,
								PayloadType: payloadType,
								Payload:     payload,
							},
							{
								ID:          requestId2,
								PayloadType: payloadType,
								Payload:     payload,
							},
						},
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Not returning failed messages")
					Expect(res.FailedTasks).To(HaveLen(0))
				})
			})

			When("An event fails to publish", func() {
				It("Return the failed event", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Push with multiple events")
					mockQueueClient.EXPECT().SendBatch(gomock.Any(), &v1.QueueSendBatchRequest{
						Queue: "test-queue",
						Tasks: []*v1.NitricTask{
							{
								Id:          "test-request-id",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
							{
								Id:          "test-request-id2",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
						},
					}).Return(&v1.QueueSendBatchResponse{
						FailedTasks: []*v1.FailedTask{
							{
								Task:    &v1.NitricTask{},
								Message: "mock failure message",
							},
						},
					}, nil)

					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{
						"hello": 123,
					}
					requestId2 := "test-request-id2"

					client := NewWithClient(mockQueueClient)
					res, err := client.SendBatch(&SendBatchOptions{
						Queue: "test-queue",
						Tasks: []*Task{
							{
								ID:          requestId,
								PayloadType: payloadType,
								Payload:     payload,
							},
							{
								ID:          requestId2,
								PayloadType: payloadType,
								Payload:     payload,
							},
						},
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the failed message")
					Expect(res.FailedTasks).To(HaveLen(1))
				})
			})

			When("An error is returned from the gRPC client", func() {
				It("Return the failed event", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Push with multiple events")
					mockQueueClient.EXPECT().SendBatch(gomock.Any(), gomock.Any()).Return(nil,
						fmt.Errorf("mock error"))

					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{}

					client := NewWithClient(mockQueueClient)
					_, err := client.SendBatch(&SendBatchOptions{
						Queue: "test-queue",
						Tasks: []*Task{
							{
								ID:          requestId,
								PayloadType: payloadType,
								Payload:     payload,
							},
						},
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})

	When("Pop", func() {
		When("The queue exists", func() {
			When("A single item is on the queue", func() {
				It("Should return the topics", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Pop")
					mockQueueClient.EXPECT().Receive(gomock.Any(), &v1.QueueReceiveRequest{
						Queue: "test-queue",
						Depth: 5,
					}).Return(&v1.QueueReceiveResponse{
						Tasks: []*v1.NitricTask{
							&v1.NitricTask{
								Id:          "test-request-id",
								PayloadType: "test-payload-type",
								LeaseId:     "test-lease-id",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
						},
					}, nil)

					client := NewWithClient(mockQueueClient)
					items, err := client.Receive(&RecieveOptions{
						Queue: "test-queue",
						Depth: 5,
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the queue item")
					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{
						"hello": float64(123),
					}
					Expect(items.Tasks).To(BeEquivalentTo([]*Task{
						{
							Payload:     payload,
							PayloadType: payloadType,
							ID:          requestId,
							LeaseID:     "test-lease-id",
						},
					}))
				})
			})

			When("A multiple item are on the queue", func() {
				It("Should return the topics", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Pop")
					mockQueueClient.EXPECT().Receive(gomock.Any(), &v1.QueueReceiveRequest{
						Queue: "test-queue",
						Depth: 5,
					}).Return(&v1.QueueReceiveResponse{
						Tasks: []*v1.NitricTask{
							&v1.NitricTask{
								Id:          "test-request-id",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
								LeaseId: "test-lease-id",
							},
							&v1.NitricTask{
								Id:          "test-request-id2",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(345),
									},
								},
								LeaseId: "test-lease-id2",
							},
						},
					}, nil)

					client := NewWithClient(mockQueueClient)
					items, err := client.Receive(&RecieveOptions{
						Queue: "test-queue",
						Depth: 5,
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the queue item")
					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{
						"hello": float64(123),
					}
					requestId2 := "test-request-id2"
					payload2 := map[string]interface{}{
						"hello": float64(345),
					}
					Expect(items.Tasks).To(BeEquivalentTo([]*Task{
						{
							Payload:     payload,
							PayloadType: payloadType,
							ID:          requestId,
							LeaseID:     "test-lease-id",
						},
						{
							Payload:     payload2,
							PayloadType: payloadType,
							ID:          requestId2,
							LeaseID:     "test-lease-id2",
						},
					}))
				})
			})

			When("No items are on the queue", func() {
				It("Should return an empty slice of topics", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Pop")
					mockQueueClient.EXPECT().Receive(gomock.Any(), &v1.QueueReceiveRequest{
						Queue: "test-queue",
						Depth: 5,
					}).Return(&v1.QueueReceiveResponse{
						Tasks: []*v1.NitricTask{},
					}, nil)

					client := NewWithClient(mockQueueClient)
					items, err := client.Receive(&RecieveOptions{
						Queue: "test-queue",
						Depth: 5,
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the queue item")

					Expect(items.Tasks).To(Equal([]*Task{}))
				})
			})

			When("The requested depth is less than 1", func() {
				It("Should set depth to 1", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Pop with a depth of 1")
					mockQueueClient.EXPECT().Receive(gomock.Any(), &v1.QueueReceiveRequest{
						Queue: "test-queue",
						Depth: 1,
					}).Return(&v1.QueueReceiveResponse{
						Tasks: []*v1.NitricTask{},
					}, nil)

					client := NewWithClient(mockQueueClient)
					// pass in a depth less than 1
					items, err := client.Receive(&RecieveOptions{
						Queue: "test-queue",
						Depth: 0,
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the queue item")
					Expect(items.Tasks).To(Equal([]*Task{}))
				})
			})
		})

		When("An error is returned from the gRPC client", func() {
			It("Should return an empty slice of topics", func() {
				mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

				By("Calling Pop")
				mockQueueClient.EXPECT().Receive(gomock.Any(), &v1.QueueReceiveRequest{
					Queue: "test-queue",
					Depth: 5,
				}).Return(nil, fmt.Errorf("mock error"))

				client := NewWithClient(mockQueueClient)
				_, err := client.Receive(&RecieveOptions{
					Queue: "test-queue",
					Depth: 5,
				})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("Complete", func() {
		When("A task can be completed", func() {
			It("Should complete successfully", func() {
				mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

				By("Calling Complete")
				mockQueueClient.EXPECT().Complete(gomock.Any(), &v1.QueueCompleteRequest{
					Queue:   "test-queue",
					LeaseId: "test-lease-id",
				}).Return(&v1.QueueCompleteResponse{}, nil)

				client := NewWithClient(mockQueueClient)
				_, err := client.Complete(&CompleteOptions{
					Queue: "test-queue",
					Task: &Task{
						LeaseID: "test-lease-id",
					},
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("A task can't be completed", func() {
			It("Should return an error", func() {
				mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

				By("Calling Complete")
				mockQueueClient.EXPECT().Complete(gomock.Any(), gomock.Any()).Return(nil,
					fmt.Errorf("mock error"))

				client := NewWithClient(mockQueueClient)
				_, err := client.Complete(&CompleteOptions{
					Queue: "test-queue",
					Task: &Task{
						LeaseID: "test-lease-id",
					},
				})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
