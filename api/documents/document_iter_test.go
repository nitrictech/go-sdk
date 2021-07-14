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
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ = Describe("DocumentIter", func() {
	ctrl := gomock.NewController(GinkgoT())
	Context("Next", func() {
		pbstr, _ := structpb.NewStruct(map[string]interface{}{
			"test": "test",
		})

		When("the stream returns a result", func() {
			strc := mock_v1.NewMockDocumentService_QueryStreamClient(ctrl)
			strc.EXPECT().Recv().Return(&v1.DocumentQueryStreamResponse{
				Document: &v1.Document{
					Content: pbstr,
				},
			}, nil)

			di := &documentIterImpl{
				str: strc,
			}

			doc, err := di.Next()

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should contain the returned document context", func() {
				Expect(doc.Content()).To(Equal(map[string]interface{}{
					"test": "test",
				}))
			})
		})

		When("the stream returns an error", func() {
			strc := mock_v1.NewMockDocumentService_QueryStreamClient(ctrl)
			strc.EXPECT().Recv().Return(nil, fmt.Errorf("mock-error"))

			di := &documentIterImpl{
				str: strc,
			}

			_, err := di.Next()

			It("should pass through the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("mock-error"))
			})
		})

	})
})
