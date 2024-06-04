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

package keyvalue

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("KeyValue API", func() {
	var (
		ctrl      *gomock.Controller
		mockKV    *mock_v1.MockKvStoreClient
		kv        *keyValueImpl
		store     Store
		storeI    *storeImpl
		ok        bool
		storeName string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockKV = mock_v1.NewMockKvStoreClient(ctrl)
		kv = &keyValueImpl{kvClient: mockKV}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Store()", func() {
		Context("Given a valid KvStoreClient", func() {
			BeforeEach(func() {
				storeName = "test-store"
				store = kv.Store(storeName)
				storeI, ok = store.(*storeImpl)
			})

			When("creating new Store instance", func() {
				It("should return an instance of storeImpl", func() {
					Expect(ok).To(BeTrue())
				})

				It("should have the provided store name", func() {
					Expect(storeI.name).To(Equal(storeName))
				})

				It("should share the KeyValue store's gRPC client", func() {
					Expect(storeI.kvClient).To(Equal(mockKV))
				})
			})
		})
	})

	Describe("New method", func() {
		When("constructing a new queue client without the membrane", func() {
			BeforeEach(func() {
				os.Setenv("NITRIC_SERVICE_DIAL_TIMEOUT", "10")
			})
			AfterEach(func() {
				os.Unsetenv("NITRIC_SERVICE_DIAL_TIMEOUT")
			})

			c, err := New()

			It("should return a nil client", func() {
				Expect(c).To(BeNil())
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		PWhen("constructing a new queue client without dial blocking", func() {
			// TODO:
		})
	})
})

// TODO: new method testing is pending
