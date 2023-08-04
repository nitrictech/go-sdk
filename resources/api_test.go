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

package resources

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/go-sdk/faas"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

var _ = Describe("api", func() {
	Context("New", func() {
		It("can register one method per routes", func() {
			m := &manager{
				workers:  map[string]Starter{},
				builders: map[string]faas.HandlerBuilder{},
			}
			a := &api{
				name:    "testApi",
				routes:  map[string]Route{},
				manager: m,
			}
			a.Get("objects/", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Post("objects/", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Get("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Put("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Patch("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Delete("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Options("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })

			Expect(m.workers["route:testApi/objects/GET"]).ToNot(BeNil())
			Expect(m.workers["route:testApi/objects/POST"]).ToNot(BeNil())
			Expect(m.workers["route:testApi/objects/:id/GET"]).ToNot(BeNil())
			Expect(m.workers["route:testApi/objects/:id/PUT"]).ToNot(BeNil())
			Expect(m.workers["route:testApi/objects/:id/DELETE"]).ToNot(BeNil())
			Expect(m.workers["route:testApi/objects/:id/OPTIONS"]).ToNot(BeNil())
		})
		It("can get api details", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockClient := mock_v1.NewMockResourceServiceClient(ctrl)
			mockConn := mock_v1.NewMockClientConnInterface(ctrl)
			m := &manager{
				workers:  map[string]Starter{},
				builders: map[string]faas.HandlerBuilder{},
				rsc:      mockClient,
				conn:     mockConn,
			}

			mockClient.EXPECT().Details(gomock.Any(), &v1.ResourceDetailsRequest{
				Resource: &v1.Resource{
					Name: "testApi",
					Type: v1.ResourceType_Api,
				},
			}).Return(&v1.ResourceDetailsResponse{
				Id:       "1234",
				Provider: "aws",
				Service:  "lambda",
				Details: &v1.ResourceDetailsResponse_Api{
					Api: &v1.ApiResourceDetails{
						Url: "example.com/aws/thing",
					},
				},
			}, nil)
			a := &api{
				name:    "testApi",
				routes:  map[string]Route{},
				manager: m,
			}
			ad, err := a.Details(context.TODO())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ad).To(Equal(&ApiDetails{
				Details: Details{
					ID:       "1234",
					Provider: "aws",
					Service:  "lambda",
				},
				URL: "example.com/aws/thing",
			}))
		})
	})

	Context("New With Security", func() {
		It("declares a new API resource with security rules", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockClient := mock_v1.NewMockResourceServiceClient(ctrl)
			mockConn := mock_v1.NewMockClientConnInterface(ctrl)
			m := &manager{
				workers:  map[string]Starter{},
				builders: map[string]faas.HandlerBuilder{},
				rsc:      mockClient,
				conn:     mockConn,
			}

			mockClient.EXPECT().Declare(context.TODO(), &v1.ResourceDeclareRequest{
				Resource: &v1.Resource{
					Type: v1.ResourceType_Api,
					Name: "testApi",
				}, Config: &v1.ResourceDeclareRequest_Api{
					Api: &v1.ApiResource{
						SecurityDefinitions: map[string]*v1.ApiSecurityDefinition{
							"jwt": {
								Definition: &v1.ApiSecurityDefinition_Jwt{
									Jwt: &v1.ApiSecurityDefinitionJwt{
										Audiences: []string{"test"},
										Issuer:    "https://test.com",
									},
								},
							},
						},
						Security: map[string]*v1.ApiScopes{
							"jwt": {
								Scopes: []string{},
							},
						},
					},
				},
			}).Times(1)

			a, err := m.newApi("testApi", WithSecurityJwtRule("jwt", JwtSecurityRule{
				Audiences: []string{"test"},
				Issuer:    "https://test.com",
			}), WithSecurity("jwt", []string{}))

			Expect(err).ShouldNot(HaveOccurred())
			Expect(a).ToNot(BeNil())

			ctrl.Finish()
		})
	})
})
