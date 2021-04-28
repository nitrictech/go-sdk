// Copyright 2021 Nitric Pty Ltd.
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

package topicclient

import (
	"fmt"

	"github.com/golang/mock/gomock"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Topicclient", func() {
	ctrl := gomock.NewController(GinkgoT())

	When("GetTopics", func() {
		When("Topics are available", func() {
			It("Should return the topics", func() {
				mockTopicClient := mock_v1.NewMockTopicClient(ctrl)

				By("Calling GetTopics")
				mockTopicClient.EXPECT().List(gomock.Any(), &v1.TopicListRequest{}).Return(&v1.TopicListResponse{
					Topics: []*v1.NitricTopic{{Name: "test-topic"}},
				}, nil)

				client := NewWithClient(nil, mockTopicClient)
				topics, err := client.GetTopics()

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning the topics")
				Expect(topics).To(Equal([]Topic{
					&NitricTopic{
						name: "test-topic",
					},
				}))
			})
		})

		When("An error is returned from the gRPC client", func() {
			It("Should return an error", func() {
				mockTopicClient := mock_v1.NewMockTopicClient(ctrl)

				By("Calling GetTopics")
				mockTopicClient.EXPECT().List(gomock.Any(), &v1.TopicListRequest{}).Return(nil, fmt.Errorf("mock error"))

				client := NewWithClient(nil, mockTopicClient)
				_, err := client.GetTopics()

				By("Not returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("Topic", func() {
		When("A topic has been instantiated", func() {
			topic := NitricTopic{
				name: "test-topic",
			}

			It("should return its name", func() {
				Expect(topic.GetName()).To(Equal("test-topic"))
			})

			It("should be printable", func() {
				Expect(topic.String()).To(Equal("test-topic"))
			})
		})
	})
})
