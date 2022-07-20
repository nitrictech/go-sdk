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

	v1 "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/faas"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("api", func() {
	Context("New", func() {
		It("can register multiple methods on routes", func() {
			m := &manager{
				blockers: map[string]Starter{},
				builders: map[string]faas.HandlerBuilder{},
			}
			a := &api{
				name:   "testApi",
				routes: map[string]Route{},
				m:      m,
			}
			a.Get("objects/", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Post("objects/", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Get("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Put("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Patch("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Delete("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })
			a.Options("objects/:id", func(hc *faas.HttpContext, hh faas.HttpHandler) (*faas.HttpContext, error) { return hc, nil })

			Expect(m.blockers["route:testApi/objects/:id"]).ToNot(BeNil())
			Expect(m.blockers["route:testApi/objects"]).ToNot(BeNil())
			Expect(m.builders["testApi/objects/:id"].String()).To(Equal("Api:testApi, path:objects/:id methods:[DELETE,GET,OPTIONS,PATCH,PUT]"))
			Expect(m.builders["testApi/objects"].String()).To(Equal("Api:testApi, path:objects/ methods:[GET,POST]"))
		})
	})

	Context("New With Security", func() {
		It("declares a new API resource with security rules", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockClient := mock_v1.NewMockResourceServiceClient(ctrl)
			mockConn := mock_v1.NewMockClientConnInterface(ctrl)
			m := &manager{
				blockers: map[string]Starter{},
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

			a, err := m.NewApi("testApi", WithSecurityJwtRule("jwt", JwtSecurityRule{
				Audiences: []string{"test"},
				Issuer:    "https://test.com",
			}), WithSecurity("jwt", []string{}))

			Expect(err).ShouldNot(HaveOccurred())
			Expect(a).ToNot(BeNil())

			ctrl.Finish()
		})
	})
})
