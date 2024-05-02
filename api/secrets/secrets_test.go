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

package secrets

// import (
// 	"os"

// 	"github.com/golang/mock/gomock"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"

// 	mock_v1 "github.com/nitrictech/go-sdk/mocks"
// )

// var _ = Describe("Secrets", func() {
// 	Context("New", func() {
// 		When("Constructing a new Secrets client with no rpc server available", func() {
// 			os.Setenv("NITRIC_SERVICE_DIAL_TIMEOUT", "10")
// 			c, err := New()

// 			It("should return a nil client", func() {
// 				Expect(c).To(BeNil())
// 			})

// 			It("should return an error", func() {
// 				Expect(err).To(HaveOccurred())
// 			})
// 		})
// 	})

// 	Context("Secret", func() {
// 		When("Retrieving a new secret reference", func() {
// 			ctrl := gomock.NewController(GinkgoT())
// 			mc := mock_v1.NewMockSecretServiceClient(ctrl)
// 			c := &secretsImpl{
// 				secretClient: mc,
// 			}

// 			It("should successfully return a correct secret reference", func() {
// 				s := c.Secret("test")

// 				By("returning s secretRefImpl")
// 				_, ok := s.(*secretRefImpl)
// 				Expect(ok).To(BeTrue())

// 				By("containing the requested name")
// 				Expect(s.Name()).To(Equal("test"))
// 			})
// 		})
// 	})
// })
