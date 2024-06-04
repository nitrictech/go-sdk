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
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/nitrictech/go-sdk/api/errors/codes"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

var _ = Describe("File", func() {
	var (
		ctrl        *gomock.Controller
		mockStorage *mock_v1.MockStorageClient
		bucket      *bucketImpl
		file        File
		bucketName  string
		fileName    string
		ctx         context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mock_v1.NewMockStorageClient(ctrl)

		bucketName = "test-bucket"
		fileName = "test-file.txt"

		bucket = &bucketImpl{
			name:          bucketName,
			storageClient: mockStorage,
		}
		file = bucket.File(fileName)

		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Name()", func() {
		It("should have the same file name as the one provided", func() {
			_fileName := file.Name()
			Expect(_fileName).To(Equal(fileName))
		})
	})

	Describe("Read()", func() {
		When("the gRPC Read operation is successful", func() {
			var fileContent []byte

			BeforeEach(func() {
				fileContent = []byte("this is dummy file content for testing")

				By("the gRPC server returning a successful response")
				mockStorage.EXPECT().Read(gomock.Any(), &v1.StorageReadRequest{
					BucketName: bucketName,
					Key:        fileName,
				}).Return(&v1.StorageReadResponse{
					Body: fileContent,
				}, nil).Times(1)
			})

			It("should return the read bytes", func() {
				fileData, err := file.Read(ctx)

				By("not returning any error")
				Expect(err).ToNot(HaveOccurred())

				By("returning the expected data in file")
				Expect(fileData).To(Equal(fileContent))
			})
		})

		When("the grpc server returns an error", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "Internal Error"

				By("the gRPC server returning an error")
				mockStorage.EXPECT().Read(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New(errorMsg),
				).Times(1)
			})

			It("should return the passed error", func() {
				fileData, err := file.Read(ctx)

				By("returning error with expected message")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())

				By("returning nil as file data")
				Expect(fileData).To(BeNil())
			})
		})
	})

	Describe("Write", func() {
		var fileData []byte

		BeforeEach(func() {
			fileData = []byte("this is dummy file content for testing")
		})

		When("the gRPC write operation is successful", func() {
			BeforeEach(func() {
				By("the gRPC server returning a successful response")
				mockStorage.EXPECT().Write(gomock.Any(), &v1.StorageWriteRequest{
					BucketName: bucketName,
					Key:        fileName,
					Body:       fileData,
				}).Return(
					&v1.StorageWriteResponse{},
					nil,
				).Times(1)
			})

			It("should not return an error", func() {
				err := file.Write(ctx, fileData)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("the grpc server returns an error", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "Internal Error"

				By("the gRPC server returning an error")
				mockStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New(errorMsg),
				).Times(1)
			})

			It("should return the passed error", func() {
				err := file.Write(ctx, fileData)

				By("returning error with expected message")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
			})
		})
	})

	Describe("Delete", func() {
		When("the delete gRPC operation is successful", func() {
			BeforeEach(func() {
				By("the gRPC server returning a successful response")
				mockStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(
					&v1.StorageDeleteResponse{},
					nil,
				).Times(1)
			})

			It("should not return an error", func() {
				err := file.Delete(ctx)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("the grpc server returns an error", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "Internal Error"

				By("the gRPC server returning an error")
				mockStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New(errorMsg),
				).Times(1)
			})

			It("should pass through the returned error", func() {
				err := file.Delete(ctx)

				By("returning error with expected message")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
			})
		})
	})

	Describe("UploadUrl()", func() {
		var expiry int

		BeforeEach(func() {
			expiry = 1 * 24 * 60 * 60 // 1 Day
		})

		When("the PreSignUrl gRPC operation is successful", func() {
			var url string

			BeforeEach(func() {
				url = "https://example.com"

				mockStorage.EXPECT().PreSignUrl(ctx, &v1.StoragePreSignUrlRequest{
					BucketName: bucketName,
					Key:        fileName,
					Operation:  v1.StoragePreSignUrlRequest_WRITE,
					Expiry:     durationpb.New(time.Duration(expiry) * time.Second),
				}).Return(&v1.StoragePreSignUrlResponse{
					Url: url,
				}, nil).Times(1)
			})

			It("should return a valid url stirng", func() {
				_url, err := file.UploadUrl(ctx, expiry)

				By("not returning any errors")
				Expect(err).ToNot(HaveOccurred())
				By("returning a valid url")
				Expect(_url).To(Equal(url))
			})
		})

		When("the grpc server returns an error", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "Internal Error"

				mockStorage.EXPECT().PreSignUrl(gomock.Any(), gomock.Any()).Return(
					&v1.StoragePreSignUrlResponse{
						Url: "",
					},
					errors.New(errorMsg),
				).Times(1)
			})

			It("should pass through the returned error", func() {
				_url, err := file.UploadUrl(ctx, expiry)

				By("returning error with expected message")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())

				By("returning empty string url")
				Expect(_url).To(Equal(""))
			})
		})
	})

	Describe("DownloadUrl()", func() {
		var expiry int

		BeforeEach(func() {
			expiry = 1 * 24 * 60 * 60 // 1 Day
		})

		When("the PreSignUrl gRPC operation is successful", func() {
			var url string

			BeforeEach(func() {
				url = "https://example.com"

				mockStorage.EXPECT().PreSignUrl(ctx, &v1.StoragePreSignUrlRequest{
					BucketName: bucketName,
					Key:        fileName,
					Operation:  v1.StoragePreSignUrlRequest_READ,
					Expiry:     durationpb.New(time.Duration(expiry) * time.Second),
				}).Return(&v1.StoragePreSignUrlResponse{
					Url: url,
				}, nil).Times(1)
			})

			It("should return a valid url stirng", func() {
				_url, err := file.DownloadUrl(ctx, expiry)

				By("not returning any errors")
				Expect(err).ToNot(HaveOccurred())
				By("returning a valid url")
				Expect(_url).To(Equal(url))
			})
		})

		When("the grpc server returns an error", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "Internal Error"

				mockStorage.EXPECT().PreSignUrl(gomock.Any(), gomock.Any()).Return(
					&v1.StoragePreSignUrlResponse{
						Url: "",
					},
					errors.New(errorMsg),
				).Times(1)
			})

			It("should pass through the returned error", func() {
				_url, err := file.DownloadUrl(ctx, expiry)

				By("returning error with expected message")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())

				By("returning empty string url")
				Expect(_url).To(Equal(""))
			})
		})
	})
})

var _ = Describe("fileImpl", func() {
	var (
		ctrl        *gomock.Controller
		mockStorage *mock_v1.MockStorageClient
		bucket      *bucketImpl
		file        File
		fileI       *fileImpl
		bucketName  string
		fileName    string
		ok          bool
		ctx         context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mock_v1.NewMockStorageClient(ctrl)

		bucketName = "test-bucket"
		fileName = "test-file.txt"

		bucket = &bucketImpl{
			name:          bucketName,
			storageClient: mockStorage,
		}
		file = bucket.File(fileName)

		By("accessing fileImpl from file")
		fileI, ok = file.(*fileImpl)
		Expect(ok).To(BeTrue())

		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("signUrl()", func() {
		var expiry int

		BeforeEach(func() {
			expiry = 1 * 24 * 60 * 60 // 1 day
		})

		When("invalid mode is provided", func() {
			It("should return an error", func() {
				_, err := fileI.signUrl(ctx, PresignUrlOptions{
					Mode:   9999, // Invalid Mode
					Expiry: expiry,
				})
				Expect(err).Should(HaveOccurred())
				Expect(strings.Contains(err.Error(), codes.InvalidArgument.String())).To(BeTrue())
			})
		})

		When("The grpc server returns an error", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "Internal Error"

				By("the gRPC server returning an error")
				mockStorage.EXPECT().PreSignUrl(gomock.Any(), &v1.StoragePreSignUrlRequest{
					BucketName: bucketName,
					Key:        fileName,
					Operation:  v1.StoragePreSignUrlRequest_READ,
					Expiry:     durationpb.New(time.Duration(expiry) * time.Second),
				}).Return(nil, errors.New(errorMsg)).Times(1)
			})

			It("should pass through the returned error", func() {
				_url, err := fileI.signUrl(ctx, PresignUrlOptions{
					Mode:   ModeRead,
					Expiry: expiry,
				})

				By("returning error with expected message")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())

				By("returning empty string url")
				Expect(_url).To(Equal(""))
			})
		})

		When("the PreSignUrl operation is successful", func() {
			var url string
			var mode Mode

			BeforeEach(func() {
				url = "http://example.com"
				mode = ModeWrite

				mockStorage.EXPECT().PreSignUrl(gomock.Any(), &v1.StoragePreSignUrlRequest{
					BucketName: bucketName,
					Key:        fileName,
					Operation:  v1.StoragePreSignUrlRequest_WRITE,
					Expiry:     durationpb.New(time.Duration(expiry) * time.Second),
				}).Return(&v1.StoragePreSignUrlResponse{
					Url: url,
				}, nil)
			})

			It("should return a success response", func() {
				_url, err := fileI.signUrl(ctx, PresignUrlOptions{Mode: mode, Expiry: expiry})

				By("no error being returned")
				Expect(err).ToNot(HaveOccurred())

				By("return expected url string")
				Expect(_url).To(Equal(url))
			})
		})
	})
})

var _ = Describe("PresignUrlOptions", func() {
	var mode Mode
	var expiry int
	var p *PresignUrlOptions

	Describe("isValid()", func() {
		When("valid mode and expiry are passed", func() {
			BeforeEach(func() {
				expiry = 1 * 24 * 60 * 60 // 1 day
				mode = ModeRead

				p = &PresignUrlOptions{
					Mode:   mode,
					Expiry: expiry,
				}
			})

			It("should not return an error", func() {
				err := p.isValid()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("invalid mode is passed", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "invalid mode"
				expiry = 1 * 24 * 60 * 60 // 1 day

				p = &PresignUrlOptions{
					Mode:   7,
					Expiry: expiry,
				}
			})

			It("should return an error", func() {
				err := p.isValid()
				By("occurance of error")
				Expect(err).To(HaveOccurred())

				By("containing appropriate error message")
				Expect(strings.Contains(
					strings.ToLower(err.Error()),
					strings.ToLower(errorMsg),
				),
				).To(BeTrue())
			})
		})

		When("invalid expiry is passed", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "invalid expiry"
				expiry = 9999 * 24 * 60 * 60 // 9999 days
				mode = ModeRead

				p = &PresignUrlOptions{
					Mode:   mode,
					Expiry: expiry,
				}
			})

			It("should return an error", func() {
				err := p.isValid()
				By("occurance of error")
				Expect(err).To(HaveOccurred())

				By("containing appropriate error message")
				Expect(strings.Contains(
					strings.ToLower(err.Error()),
					strings.ToLower(errorMsg),
				),
				).To(BeTrue())
			})
		})
	})
})
