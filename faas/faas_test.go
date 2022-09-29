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
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/mock/gomock"

	pb "github.com/nitrictech/apis/go/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("Faas", func() {
	Context("New", func() {
		When("Creating a new HandlerBuilder", func() {
			fs := New()

			It("Should be an instance of *faasClientImpl", func() {
				_, ok := fs.(*faasClientImpl)

				Expect(ok).To(BeTrue())
			})
		})
	})

	Context("faasClientImpl", func() {
		Context("Http", func() {
			When("Setting the Http Middleware", func() {
				mware := func(ctx *HttpContext, next HttpHandler) (*HttpContext, error) {
					return ctx, nil
				}

				impl := &faasClientImpl{http: map[string]HttpMiddleware{}}
				impl.Http("GET", mware)

				It("should set the private http field", func() {
					Expect(impl.Http("GET")).ToNot(BeNil())
				})

				When("Getting the Http Middleware", func() {
					mw := impl.GetHttp("GET")

					It("should return the internal http field", func() {
						Expect(reflect.ValueOf(impl.http["GET"]).Pointer()).To(Equal(reflect.ValueOf(mw).Pointer()))
					})
				})
			})
		})

		Context("Event", func() {
			When("Setting the Event Middleware", func() {
				mware := func(ctx *EventContext, next EventHandler) (*EventContext, error) {
					return ctx, nil
				}

				impl := &faasClientImpl{http: map[string]HttpMiddleware{}}
				impl.Event(mware)

				It("should set the private event field", func() {
					Expect(impl.event).ToNot(BeNil())
				})

				When("Getting the Event Middleware", func() {
					mw := impl.GetEvent()

					It("should return the internal event field", func() {
						Expect(reflect.ValueOf(impl.event).Pointer()).To(Equal(reflect.ValueOf(mw).Pointer()))
					})
				})
			})
		})

		Context("Default", func() {
			When("Setting the Default Middleware", func() {
				mware := func(ctx TriggerContext, next TriggerHandler) (TriggerContext, error) {
					return ctx, nil
				}

				impl := &faasClientImpl{http: map[string]HttpMiddleware{}}
				impl.Default(mware)

				It("should set the private trig field", func() {
					Expect(impl.trig).ToNot(BeNil())
				})

				When("Getting the Default Middleware", func() {
					mw := impl.GetDefault()

					It("should return the internal trig field", func() {
						Expect(reflect.ValueOf(impl.trig).Pointer()).To(Equal(reflect.ValueOf(mw).Pointer()))
					})
				})
			})
		})
	})

	Context("Start", func() {
		impl := &faasClientImpl{http: map[string]HttpMiddleware{}}
		When("No FaasServiceServer is available", func() {
			err := impl.Start()

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("A FaasServiceServer is available", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockClient := mock_v1.NewMockFaasServiceClient(ctrl)
			mockStream := mock_v1.NewMockFaasService_TriggerStreamClient(ctrl)
			When("no valid handlers are provided", func() {
				err := impl.startWithClient(mockClient)

				It("should return an error", func() {
					Expect(err).Should(HaveOccurred())
				})
			})

			When("a valid handler is provided", func() {
				impl.Http("GET", func(ctx *HttpContext, next HttpHandler) (*HttpContext, error) {
					return ctx, nil
				})
				impl.Http("POST", func(ctx *HttpContext, next HttpHandler) (*HttpContext, error) {
					return ctx, nil
				})
				impl.WithApiWorkerOpts(ApiWorkerOptions{
					ApiName: "test",
					Path:    "apples",
					Security: map[string][]string{
						"x": {"y"},
					},
					SecurityDisabled: false,
				})

				It("should start the faas loop", func() {
					By("Opening a stream with the Faas server")
					mockClient.EXPECT().TriggerStream(gomock.Any()).Return(mockStream, nil)

					By("Sending an InitRequest")
					mockStream.EXPECT().Send(&pb.ClientMessage{
						Content: &pb.ClientMessage_InitRequest{
							InitRequest: &pb.InitRequest{
								Worker: &pb.InitRequest_Api{
									Api: &pb.ApiWorker{
										Api:     "test",
										Path:    "apples",
										Methods: []string{"GET", "POST"},
										Options: &pb.ApiWorkerOptions{
											SecurityDisabled: false,
											Security: map[string]*pb.ApiWorkerScopes{
												"x": {
													Scopes: []string{"y"},
												},
											},
										},
									},
								},
							},
						},
					}).Return(nil)

					By("The stream closing on first message")
					mockStream.EXPECT().Recv().Return(nil, io.EOF)

					err := impl.startWithClient(mockClient)

					By("Returning the stream close error")
					Expect(err).Should(HaveOccurred())

					// assert prior exprects were called
					ctrl.Finish()
				})
			})
		})
	})
})
