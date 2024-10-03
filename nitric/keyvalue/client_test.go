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
	"context"
	"errors"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
	"github.com/nitrictech/protoutils"
)

var _ = Describe("KeyValue Store API", func() {
	var (
		ctrl      *gomock.Controller
		mockKV    *mock_v1.MockKvStoreClient
		kv        *KvStoreClient
		store     KvStoreClientIface
		storeName string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockKV = mock_v1.NewMockKvStoreClient(ctrl)
		storeName = "test-store"
		kv = &KvStoreClient{name: storeName, kvClient: mockKV}
		store = kv
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Having a valid store", func() {
		When("requesting a store's name", func() {
			It("should have the provided store name", func() {
				Expect(store.Name()).To(Equal(storeName))
			})
		})

		Describe("Get", func() {
			var key string
			var expectedValue map[string]interface{}

			BeforeEach(func() {
				key = "test-key"
				expectedValue = map[string]interface{}{"data": "value"}
			})

			When("the key exists", func() {
				BeforeEach(func() {
					contentStruct, _ := protoutils.NewStruct(expectedValue)
					mockKV.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(&v1.KvStoreGetValueResponse{
						Value: &v1.Value{
							Ref: &v1.ValueRef{
								Store: storeName,
								Key:   key,
							},
							Content: contentStruct,
						},
					}, nil).Times(1)
				})

				It("should return the correct value", func() {
					value, err := store.Get(context.Background(), key)
					Expect(err).NotTo(HaveOccurred())
					Expect(value).To(Equal(expectedValue))
				})
			})

			When("the key does not exists", func() {
				BeforeEach(func() {
					mockKV.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(&v1.KvStoreGetValueResponse{
						Value: nil,
					}, nil).Times(1)
				})

				It("should return an error", func() {
					_, err := store.Get(context.Background(), key)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("Set", func() {
			var key string
			var valueToSet map[string]interface{}

			BeforeEach(func() {
				key = "test-key"
				valueToSet = map[string]interface{}{"data": "value"}
			})

			When("the operation is successful", func() {
				BeforeEach(func() {
					mockKV.EXPECT().SetValue(gomock.Any(), gomock.Any()).Return(
						&v1.KvStoreSetValueResponse{},
						nil,
					).Times(1)
				})

				It("should successfully set the value", func() {
					err := store.Set(context.Background(), key, valueToSet)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("the operation fails", func() {
				var errorMsg string
				BeforeEach(func() {
					errorMsg = "Internal Error"
					mockKV.EXPECT().SetValue(gomock.Any(), gomock.Any()).Return(
						nil,
						errors.New(errorMsg),
					).Times(1)
				})

				It("should return an error", func() {
					err := store.Set(context.Background(), key, valueToSet)
					Expect(err).To(HaveOccurred())
					Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
				})
			})
		})

		Describe("Delete", func() {
			var key string

			BeforeEach(func() {
				key = "test-key"
			})

			When("the operation is successful", func() {
				BeforeEach(func() {
					mockKV.EXPECT().DeleteKey(gomock.Any(), gomock.Any()).Return(
						&v1.KvStoreDeleteKeyResponse{},
						nil,
					).Times(1)
				})

				It("should successfully set the value", func() {
					err := store.Delete(context.Background(), key)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("the GRPC operation fails", func() {
				var errorMsg string

				BeforeEach(func() {
					errorMsg = "Internal Error"
					mockKV.EXPECT().DeleteKey(gomock.Any(), gomock.Any()).Return(
						nil,
						errors.New(errorMsg),
					).Times(1)
				})

				It("should return an error", func() {
					err := store.Delete(context.Background(), key)
					Expect(err).To(HaveOccurred())
					Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
				})
			})
		})

		Describe("Keys", func() {
			When("the operation is successful", func() {
				var expectedKey string

				BeforeEach(func() {
					expectedKey = "key1"
					mockStream := mock_v1.NewMockKvStore_ScanKeysClient(ctrl)
					mockKV.EXPECT().ScanKeys(gomock.Any(), gomock.Any()).Return(mockStream, nil).Times(1)
					mockStream.EXPECT().Recv().Return(&v1.KvStoreScanKeysResponse{Key: expectedKey}, nil).AnyTimes()
				})

				It("should return a stream of keys", func() {
					stream, err := store.Keys(context.Background())
					Expect(err).ToNot(HaveOccurred())
					key, err := stream.Recv()
					Expect(err).ToNot(HaveOccurred())
					Expect(key).To(Equal(expectedKey))
				})
			})

			When("the operation fails", func() {
				var errorMsg string
				BeforeEach(func() {
					errorMsg = "Internal Error"
					mockStream := mock_v1.NewMockKvStore_ScanKeysClient(ctrl)
					mockKV.EXPECT().ScanKeys(gomock.Any(), gomock.Any()).Return(mockStream, nil).Times(1)
					mockStream.EXPECT().Recv().Return(nil, errors.New(errorMsg)).Times(1)
				})

				It("should return an error", func() {
					stream, err := store.Keys(context.Background())
					Expect(err).ToNot(HaveOccurred())
					_, err = stream.Recv()
					Expect(err).To(HaveOccurred())
					Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
				})
			})
		})
	})
})
