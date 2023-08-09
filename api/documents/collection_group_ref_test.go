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
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("CollectionGroupRef", func() {
	ctrl := gomock.NewController(GinkgoT())
	mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

	mcg := &collectionGroupRefImpl{
		documentClient: mdc,
		name:           "test-mock",
		parent: &collectionGroupRefImpl{
			documentClient: mdc,
			name:           "test-mock-parent",
		},
	}

	Context("Name", func() {
		When("retrieving the name if a collection group ref", func() {
			name := mcg.Name()

			It("should return the internal name field", func() {
				Expect(name).To(Equal(mcg.name))
			})
		})
	})

	Context("Query", func() {
		When("constructing a Query builder from a Collection group reference", func() {
			q := mcg.Query()

			qi, ok := q.(*queryImpl)

			It("should return a queryImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should contain the sub collection references", func() {
				Expect(qi.col.Name()).To(Equal("test-mock"))
			})

			It("should contain the parent collection reference", func() {
				Expect(qi.col.Parent().Parent().Name()).To(Equal("test-mock-parent"))
			})
		})
	})

	Context("Parent", func() {
		When("retrieving the parent of a collection group reference", func() {
			p := mcg.Parent()

			cgi, ok := p.(*collectionGroupRefImpl)

			It("should return a collectionGroupRefImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the name of the parent collection", func() {
				Expect(cgi.name).To(Equal("test-mock-parent"))
			})
		})
	})

	Context("toColRef", func() {
		When("converting a CollectionGroupReference to a CollectionReference", func() {
			cr := mcg.toColRef()

			cri, ok := cr.(*collectionRefImpl)

			It("should return a collectionRefImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the collection groups name", func() {
				Expect(cri.name).To(Equal(mcg.name))
			})

			It("should have a parent document with a blank key", func() {
				Expect(cri.Parent().Id()).To(Equal(""))
			})

			It("the parent document should have the given parent collection", func() {
				Expect(cri.Parent().Parent().Name()).To(Equal("test-mock-parent"))
			})
		})
	})

	Context("fromColRef", func() {
		When("converting from a CollectionReference to a CollectionGroupReference", func() {
			cgr := fromColRef(&collectionRefImpl{
				name: "test-mock",
				parentDocument: &documentRefImpl{
					documentClient: mdc,
					col: &collectionRefImpl{
						name:           "test-mock-parent",
						documentClient: mdc,
					},
					id: "test",
				},
			}, mdc)

			It("should have the name of the base collection", func() {
				Expect(cgr.name).To(Equal("test-mock"))
			})

			It("should have the name of the collection parent", func() {
				Expect(cgr.Parent().Name()).To(Equal("test-mock-parent"))
			})
		})
	})
})
