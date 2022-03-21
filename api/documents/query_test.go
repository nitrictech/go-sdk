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

	v1 "github.com/nitrictech/apis/go/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	"github.com/nitrictech/protoutils"
)

var _ = Describe("Query", func() {
	ctrl := gomock.NewController(GinkgoT())

	Context("Query", func() {
		Context("Where", func() {
			When("adding a where clause to a query", func() {
				q := &queryImpl{
					exps: make([]*queryExpression, 0),
				}

				r := q.Where(
					Condition("test").Eq(StringValue("test")),
				)

				It("should return a reference to the original query", func() {
					Expect(r).To(Equal(q))
				})

				It("should append the expression to exps", func() {
					Expect(q.exps).To(HaveLen(1))
					e := q.exps[0]
					Expect(e.field).To(Equal("test"))
					Expect(e.op).To(Equal(queryOp_EQ))
					Expect(e.val).To(Equal(StringValue("test")))
				})
			})
		})

		Context("Limit", func() {
			When("adding a limit clause to a query", func() {
				q := &queryImpl{}

				r := q.Limit(10)

				It("should return a reference to the original query", func() {
					Expect(r).To(Equal(q))
				})

				It("should set the query limit", func() {
					Expect(q.limit).To(Equal(10))
				})
			})
		})

		Context("FromPagingToken", func() {
			When("adding a paging token to scan from", func() {
				q := &queryImpl{}

				r := q.FromPagingToken(map[string]interface{}{
					"test": "test",
				})

				It("should return a reference to the original query", func() {
					Expect(r).To(Equal(q))
				})

				It("should set the paging token", func() {
					Expect(q.pt).To(Equal(map[string]interface{}{
						"test": "test",
					}))
				})
			})
		})

		Context("Fetch", func() {
			When("fetching with valid options", func() {
				When("the gRPC server returns an error", func() {
					mdc := mock_v1.NewMockDocumentServiceClient(ctrl)
					mdc.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unimplemented, "mock-error"))

					q := newQuery(&collectionRefImpl{
						name: "test",
						dc:   mdc,
					}, mdc)

					q.Limit(100)
					q.FromPagingToken(map[string]string{
						"test": "test",
					})
					q.Where(Condition("test").Eq(StringValue("test")))

					_, err := q.Fetch()

					It("should unwrap the gRPC error", func() {
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Unimplemented: mock-error: \n rpc error: code = Unimplemented desc = mock-error"))
					})
				})

				When("the gRPC server returns a successful response", func() {
					sv, _ := protoutils.NewStruct(map[string]interface{}{
						"test": "test",
					})
					mdc := mock_v1.NewMockDocumentServiceClient(ctrl)
					mdc.EXPECT().Query(gomock.Any(), gomock.Any()).Return(&v1.DocumentQueryResponse{
						Documents: []*v1.Document{
							{
								Key: &v1.Key{
									Collection: &v1.Collection{
										Name: "test",
									},
									Id: "test",
								},
								Content: sv,
							},
						},
					}, nil)

					q := newQuery(&collectionRefImpl{
						name: "test",
						dc:   mdc,
					}, mdc)

					q.Limit(100)
					q.FromPagingToken(map[string]string{
						"test": "test",
					})
					q.Where(
						Condition("test").Eq(StringValue("test")),
					)

					r, err := q.Fetch()

					It("should return documents", func() {
						By("not returning an error")
						Expect(err).ToNot(HaveOccurred())

						By("having the correct number of documents")
						Expect(r.Documents).To(HaveLen(1))

						By("documents containing the expected data")
						Expect(r.Documents[0].Content()).To(Equal(map[string]interface{}{
							"test": "test",
						}))
					})
				})
			})

			When("providing an invalid paging token", func() {
				mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

				q := newQuery(&collectionRefImpl{
					name: "test",
					dc:   mdc,
				}, mdc)

				q.FromPagingToken("blah")

				_, err := q.Fetch()

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid Argument: Query.Fetch: Paging Token invalid"))
				})
			})

			When("providing an invalid expression", func() {
				mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

				q := newQuery(&collectionRefImpl{
					name: "test",
					dc:   mdc,
				}, mdc)

				q.Where(&queryExpression{})

				_, err := q.Fetch()

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("Stream", func() {
			When("fetching with valid options", func() {
				When("the gRPC server returns an error", func() {
					mdc := mock_v1.NewMockDocumentServiceClient(ctrl)
					mdc.EXPECT().QueryStream(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock-error"))

					q := newQuery(&collectionRefImpl{
						name: "test",
						dc:   mdc,
					}, mdc)

					q.Limit(100)
					q.FromPagingToken(map[string]string{
						"test": "test",
					})
					q.Where(Condition("test").Eq(StringValue("test")))

					_, err := q.Stream()

					It("should unwrap the gRPC error", func() {
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Unknown: error from grpc library: \n mock-error"))
					})
				})

				When("the gRPC server returns a successful response", func() {
					mdc := mock_v1.NewMockDocumentServiceClient(ctrl)
					strc := mock_v1.NewMockDocumentService_QueryStreamClient(ctrl)
					mdc.EXPECT().QueryStream(gomock.Any(), gomock.Any()).Return(strc, nil)

					q := newQuery(&collectionRefImpl{
						name: "test",
						dc:   mdc,
					}, mdc)

					q.Limit(100)
					q.FromPagingToken(map[string]string{
						"test": "test",
					})
					q.Where(
						Condition("test").Eq(StringValue("test")),
					)

					r, err := q.Stream()

					It("should not return an error", func() {
						Expect(err).ToNot(HaveOccurred())
					})

					iter, ok := r.(*documentIterImpl)

					It("should have the returned a document iterator", func() {
						Expect(ok).To(BeTrue())
					})

					It("should have a reference to the returned stream client", func() {
						Expect(iter.str).To(Equal(strc))
					})
				})
			})

			When("providing an invalid expression", func() {
				mdc := mock_v1.NewMockDocumentServiceClient(ctrl)

				q := newQuery(&collectionRefImpl{
					name: "test",
					dc:   mdc,
				}, mdc)

				q.Where(&queryExpression{})

				_, err := q.Stream()

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
