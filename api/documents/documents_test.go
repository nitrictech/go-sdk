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

package documents

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("Documents", func() {
	ctrl := gomock.NewController(GinkgoT())

	Context("New", func() {
		When("constructing a new client without the membrane present", func() {
			os.Setenv("NITRIC_SERVICE_DIAL_TIMEOUT", "10")
			c, err := New()

			It("should return a nil client", func() {
				Expect(c).To(BeNil())
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("Collection", func() {
		When("Creating a collection reference", func() {
			mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

			dc := &documentsImpl{
				dc: mdc,
			}

			collection := dc.Collection("test")

			ci, ok := collection.(*collectionRefImpl)

			It("should return a collectionRefImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provided collection name", func() {
				Expect(ci.name).To(Equal("test"))
			})

			It("should share a documents client with the documents client", func() {
				Expect(ci.dc).To(Equal(dc.dc))
			})

			It("should have a nil parentDocument", func() {
				Expect(ci.parentDocument).To(BeNil())
			})
		})
	})
})
