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

package topics

import (
	"context"
	"errors"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"github.com/nitrictech/protoutils"
)

var _ = Describe("File", func() {
	var (
		ctrl     	*gomock.Controller
		mockTopic   *mock_v1.MockTopicsClient
		t 			 Topic
		topicName 	 string
		ctx 		 context.Context
	)

	BeforeEach(func ()  {
		ctrl = gomock.NewController(GinkgoT())
		mockTopic = mock_v1.NewMockTopicsClient(ctrl)

		topicName = "test-topic"
		t = &topicImpl{
			name: topicName,
			topicClient: mockTopic,
		}

		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Name()", func() {
		It("should have the same topic name as the one provided", func ()  {
			_topicName := t.Name()
			Expect(_topicName).To(Equal(topicName))
		})
	})

	Describe("Publish()", func() {
		var messageToBePublished map[string]interface{}

		BeforeEach(func ()  {
			messageToBePublished = map[string]interface{}{
				"data": "hello world",
			}
		})

		When("the gRPC Read operation is successful", func() {
			BeforeEach(func ()  {
				payloadStruct, err := protoutils.NewStruct(messageToBePublished)
				Expect(err).ToNot(HaveOccurred())

				mockTopic.EXPECT().Publish(gomock.Any(), &v1.TopicPublishRequest{
					TopicName: topicName,
					Message: &v1.TopicMessage{
						Content: &v1.TopicMessage_StructPayload{
							StructPayload: payloadStruct,
						},
					},
				}).Return(
					&v1.TopicPublishResponse{},
				nil).Times(1)
			})

			It("should not return error", func ()  {
				err := t.Publish(ctx, messageToBePublished)

				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("the grpc server returns an error", func() {
			var errorMsg string

			BeforeEach(func ()  {
				errorMsg = "Internal Error"

				By("the gRPC server returning an error")
				mockTopic.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New(errorMsg),
				).Times(1)
			})

			It("should return the passed error", func() {
				err := t.Publish(ctx, messageToBePublished)

				By("returning error with expected message")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
			})
		})
	})
})