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

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

var _ = Describe("Bucket", func() {
	var (
		ctrl     *gomock.Controller
		mockStorage   *mock_v1.MockStorageClient
		bucket       *bucketImpl
		bucketName string
		fileName string
		ctx 		context.Context
	)

	BeforeEach(func ()  {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mock_v1.NewMockStorageClient(ctrl)

		bucketName = "test-bucket"
		fileName = "test-file.txt"

		bucket = &bucketImpl{
			name:          bucketName,
			storageClient: mockStorage,
		}

		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("File()", func() {
		When("creating a new File reference", func() {
			var(
				file    File
				fileI   *fileImpl
				ok       bool
			)

			BeforeEach(func ()  {
				file = bucket.File(fileName)
				fileI, ok = file.(*fileImpl)		
			})

			It("it should succesfully return File Instance", func ()  {
				By("returning an fileImpl instance")
				Expect(ok).To(BeTrue())
			
				By("having the provided file name")
				Expect(fileI.key).To(Equal(fileName))

				By("sharing the bucket name")
				Expect(fileI.bucket).To(Equal(bucketName))

				By("sharing the Bucket's gRPC client")
				Expect(fileI.storageClient).To(Equal(mockStorage))
			})
		})
	})

	Describe("Files()", func() {
		When("the gRPC opreation of ListBlobs fails", func() {
			var errorMsg string

			BeforeEach(func ()  {
				errorMsg = "Internal Error"
				
				By("the nitric membrane returning an error")
				mockStorage.EXPECT().ListBlobs(gomock.Any(), &v1.StorageListBlobsRequest{
					BucketName: bucketName,
				}).Times(1).Return(nil, errors.New(errorMsg))
			})

			It("should return an error", func() {
				By("calling Files() on the bucket reference")
				files, err := bucket.Files(ctx)

				By("receiving an error with same error message")
				Expect(err).Should(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())

				By("receiving nil files")
				Expect(files).To(BeNil())
			})
		})

		When("the gRPC opreation of ListBlobs succeeds", func() {
			var files []File

			BeforeEach(func ()  {
				files = []File{
					bucket.File("file-1.txt"),
					bucket.File("file-2.txt"),
				}

				blobs := make([]*v1.Blob, 0, len(files))
				for _, file := range files {
					fileI, ok := file.(*fileImpl) 
					Expect(ok).To(BeTrue())
					blobs = append(blobs, &v1.Blob{
						Key: fileI.key,
					})
				}

				By("the bucket not being empty")
				mockStorage.EXPECT().ListBlobs(gomock.Any(), &v1.StorageListBlobsRequest{
					BucketName: bucketName,
				}).Return(&v1.StorageListBlobsResponse{
					Blobs: blobs,
				}, nil).Times(1)
			})

			It("should list the files in the bucket", func() {
				By("bucket.Files() being called")
				_files, err := bucket.Files(ctx)

				By("not returning an error")
				Expect(err).ToNot(HaveOccurred())

				By("returning the files")
				Expect(_files).To(HaveExactElements(files))
			})
		})
	})

	Describe("Name()", func() {
		It("should have the same name as the one provided", func ()  {
			_bucketName := bucket.Name()
			Expect(_bucketName).To(Equal(bucketName))
		})
	})
})
