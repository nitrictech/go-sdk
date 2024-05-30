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
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("Topics API", func() {
	var (
		ctrl     	 *gomock.Controller
		mockTopics   *mock_v1.MockTopicsClient
		ts		 	 Topics
	)
	
	BeforeEach(func ()  {
		ctrl = gomock.NewController(GinkgoT())
		mockTopics = mock_v1.NewMockTopicsClient(ctrl)
		
		ts = &topicsImpl{
			topicClient: mockTopics,
		}
	})
	
	AfterEach(func() {
		ctrl.Finish()
	})
	
	Describe("Topic()", func() {
		var topicName string
		var topicI *topicImpl
		var ok bool
		
		When("creating a new Topic reference", func() {
			BeforeEach(func ()  {
				topicName = "test-topic"
				topic := ts.Topic(topicName)
				topicI, ok = topic.(*topicImpl)
			})
			
			It("should return a topicImpl instance", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provied topic name", func() {
				Expect(topicI.name).To(Equal(topicName))
			})

			It("should share the storage clients gRPC client", func() {
				Expect(topicI.topicClient).To(Equal(mockTopics))
			})
		})
	})

	Describe("New()", func() {
		Context("constructing a new topics client", func() {
			When("the gRPC connection is unavailable", func() {
				BeforeEach(func() {
					os.Setenv("NITRIC_SERVICE_DIAL_TIMEOUT", "10")
				})
				AfterEach(func() {
					os.Unsetenv("NITRIC_SERVICE_DIAL_TIMEOUT")
				})
				
				ts, err := New()
				
				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					
					By("not returning a topics client")
					Expect(ts).To(BeNil())
				})
			})

			PWhen("constructing a new topics client without dial blocking", func() {
				// TODO: Mock an available server to connect to
			})
		})
	})
})
