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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/nitrictech/apis/go/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("Object", func() {

	Context("Read", func() {
		When("The grpc server returns an error", func() {

			It("should pass through the returned error", func() {
				ctrl := gomock.NewController(GinkgoT())

				mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
					sc:     mockStorage,
				}

				By("the gRPC server returning an error")
				mockStorage.EXPECT().Read(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock error"))

				_, err := obj.Read()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unknown: mock error"))

				ctrl.Finish()
			})
		})

		When("The read is successful", func() {
			It("should return the read bytes", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
					sc:     mockStorage,
				}

				By("the gRPC server returning a successful response")
				mockStorage.EXPECT().Read(gomock.Any(), gomock.Any()).Return(&v1.StorageReadResponse{
					Body: []byte("test"),
				}, nil)

				b, _ := obj.Read()
				Expect(b).To(Equal([]byte("test")))

				ctrl.Finish()
			})
		})
	})

	Context("Write", func() {
		When("The grpc server returns an error", func() {
			It("should pass through the returned error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
					sc:     mockStorage,
				}

				By("the gRPC server returning an error")
				mockStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock error"))

				err := obj.Write([]byte("test"))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unknown: mock error"))

				ctrl.Finish()
			})
		})

		When("The write is successful", func() {
			It("should not return an error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
					sc:     mockStorage,
				}

				By("the gRPC server returning a successful response")
				mockStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Return(&v1.StorageWriteResponse{}, nil)

				err := obj.Write([]byte("test"))
				Expect(err).ToNot(HaveOccurred())

				ctrl.Finish()
			})
		})
	})

	Context("Delete", func() {
		When("The grpc server returns an error", func() {
			It("should pass through the returned error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
					sc:     mockStorage,
				}

				By("the gRPC server returning an error")
				mockStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock error"))
				err := obj.Delete()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unknown: mock error"))

				ctrl.Finish()
			})
		})

		When("The delete is successful", func() {
			It("should not return an error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
					sc:     mockStorage,
				}

				By("the gRPC server returning a successful response")
				mockStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&v1.StorageDeleteResponse{}, nil)

				err := obj.Delete()
				Expect(err).ToNot(HaveOccurred())
				ctrl.Finish()
			})
		})
	})

	Context("SignUrl", func() {
		When("Invalid mode is provided", func() {
			It("should return an error", func() {
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
				}

				_, err := obj.SignUrl(PresignUrlOptions{
					Mode: 7,
				})
				Expect(err).Should(HaveOccurred())

				Expect(err.Error()).To(Equal("Invalid Argument: invalid options: \n invalid mode: 7"))
			})
		})

		When("The grpc server returns an error", func() {
			It("should pass through the returned error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
					sc:     mockStorage,
				}

				By("the gRPC server returning an error")
				mockStorage.EXPECT().PreSignUrl(gomock.Any(), &v1.StoragePreSignUrlRequest{
					BucketName: "test-bucket",
					Key:        "test-object",
					Operation:  v1.StoragePreSignUrlRequest_READ,
				}).Return(nil, fmt.Errorf("mock error"))

				_, err := obj.SignUrl(PresignUrlOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unknown: mock error"))

				ctrl.Finish()
			})
		})

		When("SignUrl is successful", func() {
			It("should not return an error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
				obj := &fileImpl{
					bucket: "test-bucket",
					key:    "test-object",
					sc:     mockStorage,
				}

				By("the gRPC server returning a successful response")
				mockStorage.EXPECT().PreSignUrl(gomock.Any(), &v1.StoragePreSignUrlRequest{
					BucketName: "test-bucket",
					Key:        "test-object",
					Operation:  v1.StoragePreSignUrlRequest_WRITE,
				}).Return(&v1.StoragePreSignUrlResponse{
					Url: "http://example.com",
				}, nil)

				url, err := obj.SignUrl(PresignUrlOptions{Mode: ModeWrite})
				Expect(err).ToNot(HaveOccurred())

				Expect(url).To(Equal("http://example.com"))

				ctrl.Finish()
			})
		})
	})
})
