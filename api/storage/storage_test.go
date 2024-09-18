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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("Storage API", func() {
	var (
		ctrl        *gomock.Controller
		mockStorage *mock_v1.MockStorageClient
		s           Storage
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mock_v1.NewMockStorageClient(ctrl)

		s = &storageImpl{
			storageClient: mockStorage,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Bucket()", func() {
		var bucketName string
		var bucketI *bucketImpl
		var ok bool

		When("creating a new Bucket reference", func() {
			BeforeEach(func() {
				bucketName = "test-bucket"
				bucket := s.Bucket(bucketName)
				bucketI, ok = bucket.(*bucketImpl)
			})

			It("should return a bucketImpl instance", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provied bucket name", func() {
				Expect(bucketI.name).To(Equal(bucketName))
			})

			It("should share the storage clients gRPC client", func() {
				Expect(bucketI.storageClient).To(Equal(mockStorage))
			})
		})
	})
})
