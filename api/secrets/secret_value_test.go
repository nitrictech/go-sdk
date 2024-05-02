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
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )

// var _ = Describe("secretValueImpl", func() {
// 	Context("Ref", func() {
// 		When("retrieving the ref of a secretValueImpl", func() {
// 			svi := &secretValueImpl{
// 				version: &secretVersionRefImpl{
// 					secret: &secretRefImpl{
// 						name: "test",
// 					},
// 					version: "test",
// 				},
// 				val: []byte("test"),
// 			}

// 			It("should return it's internal version field", func() {
// 				Expect(svi.Version()).To(Equal(svi.version))
// 			})
// 		})
// 	})

// 	Context("AsBytes", func() {
// 		When("retrieving secret value as bytes", func() {
// 			svi := &secretValueImpl{
// 				val: []byte("test"),
// 			}

// 			It("should return it's internal val field", func() {
// 				Expect(svi.AsBytes()).To(Equal([]byte("test")))
// 			})
// 		})
// 	})

// 	Context("AsString", func() {
// 		When("retrieving secret value as a string", func() {
// 			svi := &secretValueImpl{
// 				val: []byte("test"),
// 			}

// 			It("should return it's internal val field", func() {
// 				Expect(svi.AsString()).To(Equal("test"))
// 			})
// 		})
// 	})
// })
