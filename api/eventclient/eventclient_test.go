package eventclient

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ = Describe("Eventclient", func() {
	ctrl := gomock.NewController(GinkgoT())

	When("Publish", func() {
		When("The topic exists", func() {
			It("Should publish the event", func() {
				mockEventClient := mock_v1.NewMockEventClient(ctrl)

				By("Calling GetTopics")

				payload := map[string]interface{}{
					"test": "content",
				}

				payloadStruct, _ := structpb.NewStruct(payload)

				mockEventClient.EXPECT().Publish(gomock.Any(), &v1.EventPublishRequest{
					Topic: "test-topic",
					Event: &v1.NitricEvent{
						Id:          "abc123",
						PayloadType: "test-payload-type",
						Payload:     payloadStruct,
					},
				}).Return(&v1.EventPublishResponse{}, nil)

				client := NewWithClient(mockEventClient, nil)
				topicName := "test-topic"
				payloadType := "test-payload-type"
				requestId := "abc123"
				// FIXME: This interface doesn't match the others.
				result, err := client.Publish(&PublishOptions{
					Topic: topicName,
					Event: &Event{
						Payload:     payload,
						PayloadType: payloadType,
						ID:          requestId,
					},
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning the request id")
				Expect(result.RequestID).To(Equal("abc123"))
			})

			When("No request id is specified", func() {
				It("Should publish the event", func() {
					mockEventClient := mock_v1.NewMockEventClient(ctrl)

					By("Calling GetTopics")
					payload := map[string]interface{}{
						"test": "content",
					}

					mockEventClient.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(&v1.EventPublishResponse{}, nil)

					client := NewWithClient(mockEventClient, nil)
					topicName := "test-topic"
					payloadType := "test-payload-type"
					result, err := client.Publish(&PublishOptions{
						Topic: topicName,
						Event: &Event{
							Payload:     payload,
							PayloadType: payloadType,
						},
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the request id")
					uuid.MustParse(result.RequestID)
				})
			})
		})

		When("An error is returned from the gRPC client", func() {
			It("Should return an error", func() {
				mockEventClient := mock_v1.NewMockEventClient(ctrl)

				By("Calling GetTopics")
				payload := map[string]interface{}{
					"test": "content",
				}

				mockEventClient.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock error"))

				client := NewWithClient(mockEventClient, nil)
				topicName := "test-topic"
				payloadType := "test-payload-type"
				requestID := "abc123"
				_, err := client.Publish(&PublishOptions{
					Topic: topicName,
					Event: &Event{
						Payload:     payload,
						PayloadType: payloadType,
						ID:          requestID,
					},
				})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
