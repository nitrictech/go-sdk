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
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/go-sdk/nitric/v1"
	"github.com/nitrictech/protoutils"
)

var _ = Describe("DocumentRef", func() {
	ctrl := gomock.NewController(GinkgoT())
	mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

	md := &documentRefImpl{
		dc: mdc,
		id: "test-doc",
		col: &collectionRefImpl{
			name: "test-col",
			dc:   mdc,
		},
	}

	Context("Collection", func() {
		When("Creating a depth-1 sub-collection reference", func() {
			c, _ := md.Collection("test-collection")

			ci, ok := c.(*collectionRefImpl)
			It("should return a collectionRefImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the requested collection name", func() {
				Expect(ci.name).To(Equal("test-collection"))
			})

			It("should share the documentRefs document client reference", func() {
				Expect(ci.dc).To(Equal(mdc))
			})

			It("should have the document as a parent", func() {
				Expect(ci.parentDocument).To(Equal(md))
			})
		})

		When("Creating a n-depth sub-collection reference", func() {
			mc := &collectionRefImpl{
				parentDocument: &documentRefImpl{
					dc: mdc,
					col: &collectionRefImpl{
						name: "parent-collection",
						dc:   mdc,
					},
				},
			}
			mdp := &documentRefImpl{
				dc:  mdc,
				col: mc,
				id:  "test-doc",
			}

			_, err := mdp.Collection("test-collection")

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(
					"Invalid Argument: DocumentRef.Collection: Maximum collection depth: 1 exceeded",
				))
			})
		})
	})

	Context("Delete", func() {
		When("the grpc server returns an error", func() {
			mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

			mdc.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(
				nil,
				fmt.Errorf("mock-error"),
			)

			md := &documentRefImpl{
				dc: mdc,
				id: "test-doc",
				col: &collectionRefImpl{
					name: "test-col",
					dc:   mdc,
				},
			}

			err := md.Delete()

			It("should pass through the returned error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("mock-error"))
			})
		})

		When("the grpc server returns a successful response", func() {
			mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

			mdc.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(
				&v1.DocumentDeleteResponse{},
				nil,
			)

			md := &documentRefImpl{
				dc: mdc,
				id: "test-doc",
				col: &collectionRefImpl{
					name: "test-col",
					dc:   mdc,
				},
			}

			err := md.Delete()

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("Set", func() {
		When("the grpc server returns an error", func() {
			mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

			mdc.EXPECT().Set(gomock.Any(), gomock.Any()).Return(
				nil,
				status.Error(codes.Unimplemented, "mock-error"),
			)

			md := &documentRefImpl{
				dc: mdc,
				id: "test-doc",
				col: &collectionRefImpl{
					name: "test-col",
					dc:   mdc,
				},
			}

			err := md.Set(map[string]interface{}{
				"test": "test",
			})

			It("should unwrap the returned error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unimplemented: mock-error: \n rpc error: code = Unimplemented desc = mock-error"))
			})
		})

		When("the grpc server returns a successful response", func() {
			mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

			mdc.EXPECT().Set(gomock.Any(), gomock.Any()).Return(
				&v1.DocumentSetResponse{},
				nil,
			)

			md := &documentRefImpl{
				dc: mdc,
				id: "test-doc",
				col: &collectionRefImpl{
					name: "test-col",
					dc:   mdc,
				},
			}

			err := md.Set(map[string]interface{}{
				"test": "test",
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("Get", func() {
		When("the grpc server returns an error", func() {
			mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

			mdc.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
				nil,
				fmt.Errorf("mock-error"),
			)

			md := &documentRefImpl{
				dc: mdc,
				id: "test-doc",
				col: &collectionRefImpl{
					name: "test-col",
					dc:   mdc,
				},
			}

			_, err := md.Get()

			It("should pass through the returned error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unknown: error from grpc library: \n mock-error"))
			})
		})

		When("the grpc server returns a successful response", func() {
			ms, _ := protoutils.NewStruct(map[string]interface{}{
				"test": "test",
			})
			mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

			mdc.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
				&v1.DocumentGetResponse{
					Document: &v1.Document{
						Content: ms,
					},
				},
				nil,
			)

			md := &documentRefImpl{
				dc: mdc,
				id: "test-doc",
				col: &collectionRefImpl{
					name: "test-col",
					dc:   mdc,
				},
			}

			d, _ := md.Get()

			It("should provide the returned document", func() {
				Expect(d.Content()).To(Equal(map[string]interface{}{
					"test": "test",
				}))
			})
		})
	})
})
