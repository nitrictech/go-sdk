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

package secrets

import (
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/nitrictech/apis/go/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("secretVersionRefImpl", func() {
	Context("Version", func() {
		When("retrieving the version of a secretVersionRefImpl", func() {
			svi := &secretVersionRefImpl{
				version: "test",
			}

			It("should return it's internal version field", func() {
				Expect(svi.Version()).To(Equal("test"))
			})
		})
	})

	Context("Secret", func() {
		When("retrieving the parent secret of a secretVersionRefImpl", func() {
			si := &secretRefImpl{
				name: "test",
			}
			svi := &secretVersionRefImpl{
				secret: si,
			}

			It("should return it's internal secret field", func() {
				Expect(svi.Secret()).To(Equal(si))
			})
		})
	})

	Context("Access", func() {
		When("the RPC server returns successfully", func() {
			ctrl := gomock.NewController(GinkgoT())
			mc := mock_v1.NewMockSecretServiceClient(ctrl)
			sv := &secretVersionRefImpl{
				version: "test",
				secret: &secretRefImpl{
					name: "test",
				},
				sc: mc,
			}

			It("should return the secret version content", func() {
				defer ctrl.Finish()

				By("calling the service with the correct input")
				mc.EXPECT().Access(gomock.Any(), &v1.SecretAccessRequest{
					SecretVersion: &v1.SecretVersion{
						Secret: &v1.Secret{
							Name: "test",
						},
						Version: "test",
					},
				}).Return(&v1.SecretAccessResponse{
					SecretVersion: &v1.SecretVersion{
						Secret: &v1.Secret{
							Name: "test",
						},
						Version: "test",
					},
					Value: []byte("testing"),
				}, nil).Times(1)

				svv, err := sv.Access()

				By("not returning an error")
				Expect(err).ToNot(HaveOccurred())

				svi, ok := svv.(*secretValueImpl)

				By("returning a secretValueImpl")
				Expect(ok).To(BeTrue())

				By("containing the returned SecretValueRef")
				Expect(svi.version.Version()).To(Equal("test"))

				By("Containing the returned secret content")
				Expect(svi.val).To(Equal([]byte("testing")))
			})
		})

		When("the RPC server returns an error", func() {
			ctrl := gomock.NewController(GinkgoT())
			mc := mock_v1.NewMockSecretServiceClient(ctrl)
			sv := &secretVersionRefImpl{
				version: "test",
				secret: &secretRefImpl{
					name: "test",
				},
				sc: mc,
			}

			It("should pass through the error", func() {
				defer ctrl.Finish()

				By("calling the service with the correct input")
				mc.EXPECT().Access(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock-error")).Times(1)

				c, err := sv.Access()

				By("returning an error")
				Expect(err).To(HaveOccurred())

				By("returning nil secret content")
				Expect(c).To(BeNil())
			})
		})
	})
})
