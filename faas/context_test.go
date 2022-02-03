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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/nitrictech/apis/go/nitric/v1"
)

var _ = Describe("Faas", func() {
	Context("triggerContextFromGrpcTriggerRequest", func() {
		When("Translating an invalid request", func() {
			ctx, err := triggerContextFromGrpcTriggerRequest(&v1.TriggerRequest{})

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("should return nil context", func() {
				Expect(ctx).To(BeNil())
			})
		})

		When("Translating a Http request", func() {
			ctx, err := triggerContextFromGrpcTriggerRequest(&v1.TriggerRequest{
				Data:     []byte("Hello"),
				MimeType: "text/plain",
				Context: &v1.TriggerRequest_Http{
					Http: &v1.HttpTriggerContext{
						Method: "GET",
						Path:   "/test/path",
						Headers: map[string]*v1.HeaderValue{
							"Content-Type": {
								Value: []string{"text/plain"},
							},
						},
						QueryParams: map[string]*v1.QueryValue{
							"q": {Value: []string{"my-query"}},
						},
					},
				},
			})

			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should have http context", func() {
				Expect(ctx.Http()).ToNot(BeNil())
			})

			It("should have the provided data", func() {
				Expect(ctx.Http().Request.Data()).To(BeEquivalentTo([]byte("Hello")))
			})

			It("should have the provided Method", func() {
				Expect(ctx.Http().Request.Method()).To(Equal("GET"))
			})

			It("should have the provided path", func() {
				Expect(ctx.Http().Request.Path()).To(Equal("/test/path"))
			})

			It("should have the provided headers", func() {
				Expect(ctx.Http().Request.Headers()).To(BeEquivalentTo(map[string][]string{
					"Content-Type": {"text/plain"},
				}))
			})

			It("should have the provided query params", func() {
				Expect(ctx.Http().Request.Query()).To(BeEquivalentTo(map[string][]string{
					"q": {"my-query"},
				}))
			})

			It("should have initialized extra context", func() {
				Expect(ctx.Http().Extras).To(BeEquivalentTo(map[string]interface{}{}))
			})
		})

		When("Translating a Topic request", func() {
			ctx, err := triggerContextFromGrpcTriggerRequest(&v1.TriggerRequest{
				Data:     []byte("Hello"),
				MimeType: "text/plain",
				Context: &v1.TriggerRequest_Topic{
					Topic: &v1.TopicTriggerContext{
						Topic: "test-topic",
					},
				},
			})

			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should have Event context", func() {
				Expect(ctx.Event()).ToNot(BeNil())
			})

			It("should have the provided data", func() {
				Expect(ctx.Event().Request.Data()).To(BeEquivalentTo([]byte("Hello")))
			})

			It("should have the provided topic name", func() {
				Expect(ctx.Event().Request.Topic()).To(Equal("test-topic"))
			})

			It("should have initialized extra context", func() {
				Expect(ctx.Event().Extras).To(BeEquivalentTo(map[string]interface{}{}))
			})
		})
	})

	Context("triggerContextToGrpcTriggerResponse", func() {
		When("Translating HttpContext", func() {
			resp, err := triggerContextToGrpcTriggerResponse(&triggerContextImpl{
				http: &HttpContext{
					Response: &HttpResponse{
						Status: 404,
						Headers: map[string][]string{
							"Content-Type": {"text/plain"},
						},
						Body: []byte("Not found"),
					},
				},
			})

			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should have http context", func() {
				Expect(resp.GetHttp()).ToNot(BeNil())
			})

			It("should have the provided status", func() {
				Expect(resp.GetHttp().Status).To(Equal(int32(404)))
			})

			It("should have the provided headers", func() {
				Expect(resp.GetHttp().GetHeaders()).To(BeEquivalentTo(map[string]*v1.HeaderValue{
					"Content-Type": {Value: []string{"text/plain"}},
				}))
			})

			// TODO: Deprecated, remove for v1
			It("should have the provided old headers", func() {
				Expect(resp.GetHttp().GetHeadersOld()).To(BeEquivalentTo(map[string]string{ //nolint:staticcheck
					"Content-Type": "text/plain",
				}))
			})
		})

		When("Translating EventContext", func() {
			resp, err := triggerContextToGrpcTriggerResponse(&triggerContextImpl{
				event: &EventContext{
					Response: &EventResponse{
						Success: false,
					},
				},
			})

			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should have topic context", func() {
				Expect(resp.GetTopic()).ToNot(BeNil())
			})

			It("should have the provided success status", func() {
				Expect(resp.GetTopic().GetSuccess()).To(BeFalse())
			})
		})

		When("Translating invalid context", func() {
			resp, err := triggerContextToGrpcTriggerResponse(&triggerContextImpl{})

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("should return a nil response", func() {
				Expect(resp).To(BeNil())
			})
		})
	})
})
