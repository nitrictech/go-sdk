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

package faas_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/nitrictech/go-sdk/faas"
	pb "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockFunctionSpy struct {
	loggedTriggers           []*faas.NitricTrigger
	mockResponse             *faas.NitricResponse
	mockResponseStatus       int
	mockResponseHeaders      map[string]string
	mockResponseTopicSuccess bool
	mockResponseData         []byte
}

func (m *MockFunctionSpy) reset() {
	m.loggedTriggers = make([]*faas.NitricTrigger, 0)
}

func (m *MockFunctionSpy) handler(r *faas.NitricTrigger) (*faas.NitricResponse, error) {
	if m.loggedTriggers == nil {
		m.loggedTriggers = make([]*faas.NitricTrigger, 0)
	}

	m.loggedTriggers = append(m.loggedTriggers, r)

	defaultResponse := r.DefaultResponse()

	defaultResponse.SetData(m.mockResponseData)

	if defaultResponse.GetContext().IsHttp() {
		defaultResponse.GetContext().AsHttp().Headers = m.mockResponseHeaders
		defaultResponse.GetContext().AsHttp().Status = m.mockResponseStatus
	} else if defaultResponse.GetContext().IsTopic() {
		defaultResponse.GetContext().AsTopic().Success = m.mockResponseTopicSuccess
	}

	return defaultResponse, nil
}

type MockHttpOptions struct {
	payloadType string
	body        []byte
}

var _ = Describe("Faas", func() {
	Context("Start", func() {
		mockFunction := &MockFunctionSpy{
			mockResponseData: []byte("Hello"),
			mockResponseHeaders: map[string]string{
				"Content-Type": "text/plain",
			},
			mockResponseStatus:       200,
			mockResponseTopicSuccess: true,
		}

		BeforeEach(func() {
			mockFunction.reset()
		})

		go (func() {
			faas.Start(mockFunction.handler)
		})()

		When("Function is called with a HttpTrigger", func() {
			BeforeEach(func() {
				// Create the mock faas client here
				ctrl := gomock.NewController(GinkgoT())
				mockFaasServiceClient := mock_v1.NewMockFaasServiceClient(ctrl)
				mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)
				mockStream.EXPECT().Recv().Return(
					&pb.ServerMessage{
						Id: "test",
						Content: &pb.ServerMessage_TriggerRequest{
							TriggerRequest: &pb.TriggerRequest{
								Data: []byte("test"),
								Context: &pb.TriggerRequest_Http{
									Http: &pb.HttpTriggerContext{
										Method: "POST",
										Headers: map[string]string{
											"Content-Type": "text/plain",
										},
									},
								},
							},
						},
					}, nil,
				)

				mockStream.EXPECT().Send(gomock.Any()).AnyTimes().Return(nil)

				mockStream.EXPECT().Recv().Return(
					nil, fmt.Errorf("EOF"),
				)

				// The client should be called at least once
				mockFaasServiceClient.EXPECT().TriggerStream(gomock.Any()).Return(mockStream, nil)

				errchan := make(chan error)
				go (func(errchan chan error) {
					// Use error channel for blocking here..
					err := faas.StartWithClient(mockFunction.handler, mockFaasServiceClient)
					errchan <- err
				})(errchan)

				// Wait for the stream to finish
				<-errchan
			})

			It("Should receive the correct request", func() {
				By("Receiving a single request")
				Expect(mockFunction.loggedTriggers).To(HaveLen(1))

				receivedRequest := mockFunction.loggedTriggers[0]
				receivedContext := receivedRequest.GetContext()

				By("Having the trigger data")
				Expect(receivedRequest.GetData()).To(BeEquivalentTo([]byte("test")))

				By("Recieving a HTTP Request")
				Expect(receivedContext.IsHttp()).To(BeTrue())

				By("Recieving the correct method")
				Expect(receivedContext.AsHttp().Method).To(Equal("POST"))

				By("Recieving the correct headers")
				Expect(receivedContext.AsHttp().Headers).To(BeEquivalentTo(
					map[string]string{
						"Content-Type": "text/plain",
					},
				))
			})
		})

		When("The Function is called with a TopicTrigger", func() {
			BeforeEach(func() {
				// Create the mock faas client here
				ctrl := gomock.NewController(GinkgoT())
				mockFaasServiceClient := mock_v1.NewMockFaasServiceClient(ctrl)
				mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)
				mockStream.EXPECT().Recv().Return(
					&pb.ServerMessage{
						Id: "test",
						Content: &pb.ServerMessage_TriggerRequest{
							TriggerRequest: &pb.TriggerRequest{
								Data: []byte("test"),
								Context: &pb.TriggerRequest_Topic{
									Topic: &pb.TopicTriggerContext{
										Topic: "test",
									},
								},
							},
						},
					}, nil,
				)
				// Close the stream by returning an error
				mockStream.EXPECT().Recv().Return(
					nil, fmt.Errorf("EOF"),
				)

				mockStream.EXPECT().Send(gomock.Any()).AnyTimes().Return(nil)
				// The client should be called at least once
				mockFaasServiceClient.EXPECT().TriggerStream(gomock.Any()).Return(mockStream, nil)

				errchan := make(chan error)
				go (func(errchan chan error) {
					// Use error channel for blocking here..
					err := faas.StartWithClient(mockFunction.handler, mockFaasServiceClient)
					errchan <- err
				})(errchan)

				// Wait for the stream to finish
				<-errchan
			})

			It("Should have the supplied topic", func() {
				By("Receiving a single trigger")
				Expect(mockFunction.loggedTriggers).To(HaveLen(1))

				receivedRequest := mockFunction.loggedTriggers[0]
				receivedContext := receivedRequest.GetContext()

				By("Recieving topic context")
				Expect(receivedContext.IsTopic()).To(BeTrue())

				By("Having the correct topic name")
				Expect(receivedContext.AsTopic().Topic).To(Equal("test"))

			})
		})
	})
})
