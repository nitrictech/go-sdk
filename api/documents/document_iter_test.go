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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/protoutils"
)

var _ = Describe("DocumentIter", func() {
	ctrl := gomock.NewController(GinkgoT())
	mdc := mock_v1.NewMockDocumentServiceClient(ctrl)
	Context("Next", func() {
		pbstr, _ := protoutils.NewStruct(map[string]interface{}{
			"test": "test",
		})

		When("the stream returns a result", func() {
			strc := mock_v1.NewMockDocumentService_QueryStreamClient(ctrl)
			strc.EXPECT().Recv().Return(&v1.DocumentQueryStreamResponse{
				Document: &v1.Document{
					Key: &v1.Key{
						Collection: &v1.Collection{
							Name: "test",
						},
						Id: "test",
					},
					Content: pbstr,
				},
			}, nil)

			di := &documentIterImpl{
				documentClient:       mdc,
				documentStreamClient: strc,
			}

			doc, err := di.Next()

			It("should contain the returned document context", func() {
				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				By("Containing the expected content")
				Expect(doc.Content()).To(Equal(map[string]interface{}{
					"test": "test",
				}))
			})
		})

		When("the stream returns an error", func() {
			strc := mock_v1.NewMockDocumentService_QueryStreamClient(ctrl)
			strc.EXPECT().Recv().Return(nil, status.Error(codes.Aborted, "mock-error"))

			di := &documentIterImpl{
				documentClient:       mdc,
				documentStreamClient: strc,
			}

			_, err := di.Next()

			It("should pass through the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Aborted: mock-error: \n rpc error: code = Aborted desc = mock-error"))
			})
		})
	})
})
