package queueclient

import (
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "go.nitric.io/go-sdk/interfaces/nitric/v1"
	mock_v1 "go.nitric.io/go-sdk/mocks"
	"go.nitric.io/go-sdk/v1/eventclient"
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
					mockQueueClient.EXPECT().Push(gomock.Any(), &v1.PushRequest{
						Queue: "test-queue",
						Events: []*v1.NitricEvent{
							{
								RequestId: "test-request-id",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
						},
					}).Return(&v1.PushResponse{}, nil)

					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{
						"hello": 123,
					}

					client := NewWithClient(mockQueueClient)
					res, err := client.Push("test-queue", []eventclient.Event{
						{
							RequestId: &requestId,
							PayloadType: &payloadType,
							Payload: &payload,
						},
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Not returning failed messages")
					Expect(res.failedEvents).To(HaveLen(0))
				})
			})

			When("Multiple events are pushed", func() {
				It("Should push the events", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Push with multiple events")
					mockQueueClient.EXPECT().Push(gomock.Any(), &v1.PushRequest{
						Queue: "test-queue",
						Events: []*v1.NitricEvent{
							{
								RequestId: "test-request-id",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
							{
								RequestId: "test-request-id2",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
						},
					}).Return(&v1.PushResponse{}, nil)

					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{
						"hello": 123,
					}
					requestId2 := "test-request-id2"

					client := NewWithClient(mockQueueClient)
					res, err := client.Push("test-queue", []eventclient.Event{
						{
							RequestId: &requestId,
							PayloadType: &payloadType,
							Payload: &payload,
						},
						{
							RequestId: &requestId2,
							PayloadType: &payloadType,
							Payload: &payload,
						},
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Not returning failed messages")
					Expect(res.failedEvents).To(HaveLen(0))
				})
			})

			When("An event fails to publish", func() {
				It("Return the failed event", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Push with multiple events")
					mockQueueClient.EXPECT().Push(gomock.Any(), &v1.PushRequest{
						Queue: "test-queue",
						Events: []*v1.NitricEvent{
							{
								RequestId: "test-request-id",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
							{
								RequestId: "test-request-id2",
								PayloadType: "test-payload-type",
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"hello": structpb.NewNumberValue(123),
									},
								},
							},
						},
					}).Return(&v1.PushResponse{
						FailedMessages: []*v1.FailedMessage{
							&v1.FailedMessage{
								Event: &v1.NitricEvent{},
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
					res, err := client.Push("test-queue", []eventclient.Event{
						{
							RequestId: &requestId,
							PayloadType: &payloadType,
							Payload: &payload,
						},
						{
							RequestId: &requestId2,
							PayloadType: &payloadType,
							Payload: &payload,
						},
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the failed message")
					Expect(res.failedEvents).To(HaveLen(1))
				})
			})

			When("An error is returned from the gRPC client", func() {
				It("Return the failed event", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Push with multiple events")
					mockQueueClient.EXPECT().Push(gomock.Any(), gomock.Any()).Return(nil,
						fmt.Errorf("mock error"))

					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{}

					client := NewWithClient(mockQueueClient)
					_, err := client.Push("test-queue", []eventclient.Event{
						{
							RequestId: &requestId,
							PayloadType: &payloadType,
							Payload: &payload,
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
					mockQueueClient.EXPECT().Pop(gomock.Any(), &v1.PopRequest{
						Queue: "test-queue",
						Depth: 5,
					}).Return(&v1.PopResponse{
						Items: []*v1.NitricQueueItem{
							&v1.NitricQueueItem{
								Event:   &v1.NitricEvent{
									RequestId:   "test-request-id",
									PayloadType: "test-payload-type",
									Payload:     &structpb.Struct{
										Fields: map[string]*structpb.Value{
											"hello": structpb.NewNumberValue(123),
										},
									},
								},
								LeaseId: "test-lease-id",
							},
						},
					}, nil)
	
					client := NewWithClient(mockQueueClient)
					items, err := client.Pop("test-queue", 5)
	
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
	
					By("Returning the queue item")
					requestId := "test-request-id"
					payloadType := "test-payload-type"
					payload := map[string]interface{}{
						"hello": float64(123),
					}
					Expect(items).To(Equal([]QueueItem{
						QueueItem{
							event: eventclient.Event{
								Payload:     &payload,
								PayloadType: &payloadType,
								RequestId:   &requestId,
							},
							leaseId: "test-lease-id",
							queue: "test-queue",
						},
					}))
				})
			})

			When("A multiple item are on the queue", func() {
				It("Should return the topics", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Pop")
					mockQueueClient.EXPECT().Pop(gomock.Any(), &v1.PopRequest{
						Queue: "test-queue",
						Depth: 5,
					}).Return(&v1.PopResponse{
						Items: []*v1.NitricQueueItem{
							{
								Event:   &v1.NitricEvent{
									RequestId:   "test-request-id",
									PayloadType: "test-payload-type",
									Payload:     &structpb.Struct{
										Fields: map[string]*structpb.Value{
											"hello": structpb.NewNumberValue(123),
										},
									},
								},
								LeaseId: "test-lease-id",
							},
							{
								Event:   &v1.NitricEvent{
									RequestId:   "test-request-id2",
									PayloadType: "test-payload-type",
									Payload:     &structpb.Struct{
										Fields: map[string]*structpb.Value{
											"hello": structpb.NewNumberValue(345),
										},
									},
								},
								LeaseId: "test-lease-id2",
							},
						},
					}, nil)

					client := NewWithClient(mockQueueClient)
					items, err := client.Pop("test-queue", 5)

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
					Expect(items).To(Equal([]QueueItem{
						{
							event: eventclient.Event{
								Payload:     &payload,
								PayloadType: &payloadType,
								RequestId:   &requestId,
							},
							leaseId: "test-lease-id",
							queue: "test-queue",
						},
						{
							event: eventclient.Event{
								Payload:     &payload2,
								PayloadType: &payloadType,
								RequestId:   &requestId2,
							},
							leaseId: "test-lease-id2",
							queue: "test-queue",
						},
					}))
				})
			})

			When("No items are on the queue", func() {
				It("Should return an empty slice of topics", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Pop")
					mockQueueClient.EXPECT().Pop(gomock.Any(), &v1.PopRequest{
						Queue: "test-queue",
						Depth: 5,
					}).Return(&v1.PopResponse{
						Items: []*v1.NitricQueueItem{},
					}, nil)

					client := NewWithClient(mockQueueClient)
					items, err := client.Pop("test-queue", 5)

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the queue item")

					Expect(items).To(Equal([]QueueItem{}))
				})
			})

			When("The requested depth is less than 1", func() {
				It("Should set depth to 1", func() {
					mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

					By("Calling Pop with a depth of 1")
					mockQueueClient.EXPECT().Pop(gomock.Any(), &v1.PopRequest{
						Queue: "test-queue",
						Depth: 1,
					}).Return(&v1.PopResponse{
						Items: []*v1.NitricQueueItem{},
					}, nil)

					client := NewWithClient(mockQueueClient)
					// pass in a depth less than 1
					items, err := client.Pop("test-queue", 0)

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the queue item")
					Expect(items).To(Equal([]QueueItem{}))
				})
			})
		})

		When("An error is returned from the gRPC client", func() {
			It("Should return an empty slice of topics", func() {
				mockQueueClient := mock_v1.NewMockQueueClient(ctrl)

				By("Calling Pop")
				mockQueueClient.EXPECT().Pop(gomock.Any(), &v1.PopRequest{
					Queue: "test-queue",
					Depth: 5,
				}).Return(nil, fmt.Errorf("mock error"))

				client := NewWithClient(mockQueueClient)
				_, err := client.Pop("test-queue", 5)

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
				mockQueueClient.EXPECT().Complete(gomock.Any(), &v1.CompleteRequest{
					Queue: "test-queue",
					LeaseId: "test-lease-id",
				}).Return(&v1.CompleteResponse{}, nil)

				client := NewWithClient(mockQueueClient)
				err := client.Complete(QueueItem{
					// TODO: Consider changing this to a pointer to an event in the QueueItem
					event:   eventclient.Event{}, // event not needed in Complete method
					leaseId: "test-lease-id",
					queue:   "test-queue",
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
				err := client.Complete(QueueItem{
					event:   eventclient.Event{},
					leaseId: "test-lease-id",
					queue:   "test-queue",
				})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
