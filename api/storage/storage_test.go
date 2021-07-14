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
	"os"

	"github.com/golang/mock/gomock"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage", func() {
	ctrl := gomock.NewController(GinkgoT())

	Context("New", func() {
		When("constructing a new storage client", func() {
			When("the gRPC connection is unavailable", func() {
				// Set the timeout to 10 milliseconds
				os.Setenv("NITRIC_SERVICE_DIAL_TIMEOUT", "10")
				sc, err := New()

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})

				It("should not return a storage client", func() {
					Expect(sc).To(BeNil())
				})
			})

			PWhen("the gRPC connection is available", func() {
				// TODO: Mock an available server to connect to
			})
		})
	})

	Context("Bucket", func() {
		When("creating a new Bucket reference", func() {
			mockStorage := mock_v1.NewMockStorageClient(ctrl)

			sc := &storageImpl{
				sc: mockStorage,
			}

			bucket := sc.Bucket("test-bucket")
			bucketI, ok := bucket.(*bucketImpl)

			It("should return a bucketImpl instance", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provied bucket name", func() {
				Expect(bucketI.name).To(Equal("test-bucket"))
			})

			It("should share the storage clients gRPC client", func() {
				Expect(bucketI.sc).To(Equal(mockStorage))
			})
		})
	})
})
