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

package faas

import (
	"io"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/nitrictech/apis/go/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

type handlerSpy struct {
	http  bool
	event bool
	trig  bool
}

func (h *handlerSpy) Http(ctx *HttpContext, _ HttpHandler) (*HttpContext, error) {
	h.http = true
	return ctx, nil
}

func (h *handlerSpy) Event(ctx *EventContext, _ EventHandler) (*EventContext, error) {
	h.event = true
	return ctx, nil
}

func (h *handlerSpy) Trigger(ctx TriggerContext, _ TriggerHandler) (TriggerContext, error) {
	h.trig = true
	return ctx, nil
}

func awaitFaasLoop(str v1.FaasService_TriggerStreamClient, p HandlerProvider) error {
	errChan := make(chan error)

	// Begin the loop
	go faasLoop(str, p, errChan)

	return <-errChan
}

var _ = Describe("look", func() {
	Context("faasLoop", func() {

		When("receiving an error from the stream", func() {
			It("should return an error", func() {
				spy := &handlerSpy{}
				ctrl := gomock.NewController(GinkgoT())
				mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)

				By("receiving the error from the stream")
				mockStream.EXPECT().Recv().Return(nil, io.EOF)

				err := awaitFaasLoop(mockStream, &faasClientImpl{
					trig: spy.Trigger,
				})

				By("returning the error")
				Expect(err).ToNot(BeNil())

				By("not calling the handler")
				Expect(spy.trig).To(BeFalse())

				ctrl.Finish()
			})
		})

		When("receiving a http request from the stream", func() {
			mockHttpRequest := &v1.ServerMessage{
				Id: "1234",
				Content: &v1.ServerMessage_TriggerRequest{
					TriggerRequest: &v1.TriggerRequest{
						Context: &v1.TriggerRequest_Http{
							Http: &v1.HttpTriggerContext{
								Method: "GET",
								Path:   "/",
							},
						},
					},
				},
			}
			defaultHttpResponse := &v1.ClientMessage{
				Id: "1234",
				Content: &v1.ClientMessage_TriggerResponse{
					TriggerResponse: &v1.TriggerResponse{
						Data: []byte("Success"),
						Context: &v1.TriggerResponse_Http{
							Http: &v1.HttpResponseContext{
								Status: 200,
								HeadersOld: map[string]string{
									"Content-Type": "text/plain",
								},
								Headers: map[string]*v1.HeaderValue{
									"Content-Type": {Value: []string{"text/plain"}},
								},
							},
						},
					},
				},
			}

			When("there is an available http handler", func() {
				It("should call the http handler", func() {
					spy := &handlerSpy{}
					ctrl := gomock.NewController(GinkgoT())
					mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)

					By("receiving the error from the stream")
					gomock.InOrder(
						mockStream.EXPECT().Recv().Return(mockHttpRequest, nil),
						mockStream.EXPECT().Recv().Return(nil, io.EOF),
					)

					By("receiving the default response from the loop")
					mockStream.EXPECT().Send(defaultHttpResponse).Return(nil)

					err := awaitFaasLoop(mockStream, &faasClientImpl{
						http: map[string]HttpMiddleware{"GET": spy.Http},
					})

					By("return the error")
					Expect(err).Should(HaveOccurred())

					By("calling the handler")
					Expect(spy.http).To(BeTrue())

					ctrl.Finish()
				})
			})

			When("there is an available trigger handler", func() {
				It("should call the trigger handler", func() {
					spy := &handlerSpy{}
					ctrl := gomock.NewController(GinkgoT())
					mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)

					By("receiving the error from the stream")
					gomock.InOrder(
						mockStream.EXPECT().Recv().Return(mockHttpRequest, nil),
						mockStream.EXPECT().Recv().Return(nil, io.EOF),
					)

					By("receiving the default response from the loop")
					mockStream.EXPECT().Send(defaultHttpResponse).Return(nil)

					err := awaitFaasLoop(mockStream, &faasClientImpl{
						trig: spy.Trigger,
					})

					By("return the error")
					Expect(err).Should(HaveOccurred())

					By("calling the handler")
					Expect(spy.trig).To(BeTrue())

					ctrl.Finish()
				})
			})

			When("there is no available handler", func() {
				It("should return an error", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)

					By("receiving requests from the stream")
					gomock.InOrder(
						mockStream.EXPECT().Recv().Return(mockHttpRequest, nil),
						mockStream.EXPECT().Recv().Return(nil, io.EOF),
					)

					By("receiving a http error response from the loop")
					mockStream.EXPECT().Send(&v1.ClientMessage{
						Id: "1234",
						Content: &v1.ClientMessage_TriggerResponse{
							TriggerResponse: &v1.TriggerResponse{
								Data: []byte("Internal Server Error"),
								Context: &v1.TriggerResponse_Http{
									Http: &v1.HttpResponseContext{
										Status: 500,
										HeadersOld: map[string]string{
											"Content-Type": "text/plain",
										},
										Headers: map[string]*v1.HeaderValue{
											"Content-Type": {Value: []string{"text/plain"}},
										},
									},
								},
							},
						},
					}).Return(nil)

					err := awaitFaasLoop(mockStream, &faasClientImpl{http: map[string]HttpMiddleware{}})

					By("returning the error")
					Expect(err).Should(HaveOccurred())

					ctrl.Finish()
				})
			})
		})

		When("receiving an event from the stream", func() {
			mockTopicRequest := &v1.ServerMessage{
				Id: "1234",
				Content: &v1.ServerMessage_TriggerRequest{
					TriggerRequest: &v1.TriggerRequest{
						Data: []byte("test"),
						Context: &v1.TriggerRequest_Topic{
							Topic: &v1.TopicTriggerContext{
								Topic: "test-topic",
							},
						},
					},
				},
			}
			defaultTopicResponse := &v1.ClientMessage{
				Id: "1234",
				Content: &v1.ClientMessage_TriggerResponse{
					TriggerResponse: &v1.TriggerResponse{
						Data: []byte(""),
						Context: &v1.TriggerResponse_Topic{
							Topic: &v1.TopicResponseContext{
								Success: true,
							},
						},
					},
				},
			}

			When("there is an available event handler", func() {
				It("should call the event handler", func() {
					spy := &handlerSpy{}
					ctrl := gomock.NewController(GinkgoT())
					mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)

					By("receiving the error from the stream")
					gomock.InOrder(
						mockStream.EXPECT().Recv().Return(mockTopicRequest, nil),
						mockStream.EXPECT().Recv().Return(nil, io.EOF),
					)

					By("receiving the default response from the loop")
					mockStream.EXPECT().Send(defaultTopicResponse).Return(nil)

					err := awaitFaasLoop(mockStream, &faasClientImpl{
						event: spy.Event,
					})

					By("return the error")
					Expect(err).Should(HaveOccurred())

					By("calling the handler")
					Expect(spy.event).To(BeTrue())

					ctrl.Finish()
				})
			})

			When("there is an available trigger handler", func() {
				It("should call the trigger handler", func() {
					spy := &handlerSpy{}
					ctrl := gomock.NewController(GinkgoT())
					mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)

					By("receiving the error from the stream")
					gomock.InOrder(
						mockStream.EXPECT().Recv().Return(mockTopicRequest, nil),
						mockStream.EXPECT().Recv().Return(nil, io.EOF),
					)

					By("receiving the default response from the loop")
					mockStream.EXPECT().Send(defaultTopicResponse).Return(nil)

					err := awaitFaasLoop(mockStream, &faasClientImpl{
						trig: spy.Trigger,
					})

					By("return the error")
					Expect(err).Should(HaveOccurred())

					By("calling the handler")
					Expect(spy.trig).To(BeTrue())

					ctrl.Finish()
				})
			})

			When("there is no available handler", func() {
				It("should return an error", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)

					By("receiving requests from the stream")
					gomock.InOrder(
						mockStream.EXPECT().Recv().Return(mockTopicRequest, nil),
						mockStream.EXPECT().Recv().Return(nil, io.EOF),
					)

					By("receiving a topic error response from the loop")
					mockStream.EXPECT().Send(&v1.ClientMessage{
						Id: "1234",
						Content: &v1.ClientMessage_TriggerResponse{
							TriggerResponse: &v1.TriggerResponse{
								Data: []byte(""),
								Context: &v1.TriggerResponse_Topic{
									Topic: &v1.TopicResponseContext{
										Success: false,
									},
								},
							},
						},
					}).Return(nil)

					err := awaitFaasLoop(mockStream, &faasClientImpl{http: map[string]HttpMiddleware{}})

					By("returning the error")
					Expect(err).Should(HaveOccurred())

					ctrl.Finish()
				})
			})
		})
	})
})
