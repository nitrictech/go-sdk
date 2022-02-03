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

var _ = Describe("CollectionRef", func() {
	ctrl := gomock.NewController(GinkgoT())
	mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

	mc := &collectionRefImpl{
		dc:   mdc,
		name: "test-mock",
	}

	Context("Query", func() {
		When("creating a new Query builder", func() {
			q := mc.Query()

			qi, ok := q.(*queryImpl)

			It("should be of type queryImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should share a client with the creating collection ref", func() {
				Expect(qi.dc).To(Equal(mc.dc))
			})

			It("should hold a reference to the its creating collection", func() {
				Expect(qi.col).To(Equal(mc))
			})
		})
	})

	Context("Doc", func() {
		When("creating a new Document reference", func() {
			d := mc.Doc("test")

			di, ok := d.(*documentRefImpl)

			It("should be of type documentRefImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provided document key", func() {
				Expect(di.id).To(Equal("test"))
			})

			It("should share a client with it's creating collectionRef", func() {
				Expect(di.dc).To(Equal(mc.dc))
			})

			It("should have a reference to its creating collectionRef", func() {
				Expect(di.col).To(Equal(mc))
			})
		})
	})

	Context("Collection", func() {
		When("retrieving a collection group reference from a collection", func() {
			colgroup := mc.Collection("test")

			It("should have the provided collection name", func() {
				Expect(colgroup.Name()).To(Equal("test"))
			})

			It("should have the name of the parent collection", func() {
				Expect(colgroup.Parent().Name()).To(Equal(mc.name))
			})
		})
	})

	Context("toWire", func() {
		When("translating a collection without a parent reference to wire", func() {
			wc := mc.toWire()

			It("should have the same name as the collectionRef", func() {
				Expect(wc.GetName()).To(Equal(mc.name))
			})

			It("should have no parent document", func() {
				Expect(wc.GetParent()).To(BeNil())
			})
		})

		When("translating a collection reference with a parent to wire", func() {
			mpc := &collectionRefImpl{
				parentDocument: &documentRefImpl{
					id: "test-parent",
					dc: mdc,
					col: &collectionRefImpl{
						name: "test-mock-parent",
						dc:   mdc,
					},
				},
				dc:   mdc,
				name: "test-mock-child",
			}

			wc := mpc.toWire()

			It("should have the same name as the collectionRef", func() {
				Expect(wc.GetName()).To(Equal("test-mock-child"))
			})

			It("should have a non-nil parent", func() {
				Expect(wc.GetParent()).ToNot(BeNil())
			})

			parent := wc.GetParent()

			It("parent should have the linked parent document key", func() {
				Expect(parent.GetId()).To(Equal("test-parent"))
			})

			It("parent document should have linked collection", func() {
				Expect(parent.GetCollection().GetName()).To(Equal("test-mock-parent"))
			})

			It("parent document collection should not have it's own parent", func() {
				Expect(parent.GetCollection().GetParent()).To(BeNil())
			})
		})
	})
})
