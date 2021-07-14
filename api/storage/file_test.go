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

package storage

import (
	"fmt"

	"github.com/golang/mock/gomock"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Object", func() {
	ctrl := gomock.NewController(GinkgoT())

	Context("Read", func() {
		When("The grpc server returns an error", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)
			obj := &fileImpl{
				bucket: "test-bucket",
				key:    "test-object",
				sc:     mockStorage,
			}
			mockStorage.EXPECT().Read(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock error"))

			It("should pass through the returned error", func() {
				_, err := obj.Read()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("mock error"))
			})
		})

		When("The read is successful", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)
			obj := &fileImpl{
				bucket: "test-bucket",
				key:    "test-object",
				sc:     mockStorage,
			}

			mockStorage.EXPECT().Read(gomock.Any(), gomock.Any()).Return(&v1.StorageReadResponse{
				Body: []byte("test"),
			}, nil)

			It("should return the read bytes", func() {
				b, _ := obj.Read()
				Expect(b).To(Equal([]byte("test")))
			})
		})
	})

	Context("Write", func() {
		When("The grpc server returns an error", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)
			obj := &fileImpl{
				bucket: "test-bucket",
				key:    "test-object",
				sc:     mockStorage,
			}
			mockStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock error"))

			It("should pass through the returned error", func() {
				err := obj.Write([]byte("test"))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("mock error"))
			})
		})

		When("The write is successful", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)
			obj := &fileImpl{
				bucket: "test-bucket",
				key:    "test-object",
				sc:     mockStorage,
			}
			mockStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Return(&v1.StorageWriteResponse{}, nil)

			It("should not return an error", func() {
				err := obj.Write([]byte("test"))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("Delete", func() {
		When("The grpc server returns an error", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)
			obj := &fileImpl{
				bucket: "test-bucket",
				key:    "test-object",
				sc:     mockStorage,
			}
			mockStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock error"))

			It("should pass through the returned error", func() {
				err := obj.Delete()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("mock error"))
			})
		})

		When("The delete is successful", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)
			obj := &fileImpl{
				bucket: "test-bucket",
				key:    "test-object",
				sc:     mockStorage,
			}
			mockStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&v1.StorageDeleteResponse{}, nil)

			It("should not return an error", func() {
				err := obj.Delete()
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
