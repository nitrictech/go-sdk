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
	"github.com/golang/mock/gomock"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bucket", func() {
	ctrl := gomock.NewController(GinkgoT())

	Context("File", func() {
		When("creating a new File reference", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)

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
})
