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

var _ = Describe("Bucket", func() {

	Context("File", func() {
		When("creating a new File reference", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)

			bucket := &bucketImpl{
				name: "test-bucket",
				sc:   mockStorage,
			}

			object := bucket.File("test-object")
			objectI, ok := object.(*fileImpl)

			It("should return an objectImpl instance", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provided file name", func() {
				Expect(objectI.key).To(Equal("test-object"))
			})

			It("should share the buckets name", func() {
				Expect(objectI.bucket).To(Equal("test-bucket"))
			})

			It("should share the bucket references gRPC client", func() {
				Expect(objectI.sc).To(Equal(mockStorage))
			})
		})
	})

	Context("Files", func() {
		When("failing to list files in a bucket", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)
			bucketRef := &bucketImpl{
				name: "test",
				sc:   mockStorage,
			}

			It("should return an error", func() {
				By("the nitric membrane returning an error")
				mockStorage.EXPECT().ListFiles(gomock.Any(), &v1.StorageListFilesRequest{
					BucketName: "test",
				}).Times(1).Return(nil, fmt.Errorf("mock-error"))

				By("calling Files() on the bucket reference")
				files, err := bucketRef.Files()

				By("receiving nil files")
				Expect(files).To(BeNil())

				By("receiving an error")
				Expect(err).Should(HaveOccurred())
			})
		})

		When("listing files in a bucket", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockStorage := mock_v1.NewMockStorageServiceClient(ctrl)

			bucketRef := &bucketImpl{
				name: "test-bucket",
				sc:   mockStorage,
			}

			It("should list the files in the bucket", func() {
				By("the bucket not being empty")
				mockStorage.EXPECT().ListFiles(gomock.Any(), &v1.StorageListFilesRequest{
					BucketName: "test-bucket",
				}).Times(1).Return(&v1.StorageListFilesResponse{
					Files: []*v1.File{{
						Key: "test.txt",
					}},
				}, nil)

				By("bucket.Files() being called")
				files, err := bucketRef.Files()

				By("not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("returning the files")
				Expect(files).To(HaveLen(1))
				Expect(files[0].Name()).To(Equal("test.txt"))
			})
		})
	})
})
