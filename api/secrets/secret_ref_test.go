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
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("secretRefImpl", func() {
	Context("Name", func() {
		When("retrieving the name from secretRefImpl", func() {
			s := &secretRefImpl{
				name: "test",
			}

			It("should return internal name", func() {
				Expect(s.Name()).To(Equal(s.name))
			})
		})
	})

	Context("Version", func() {
		When("retrieving a new secret version reference", func() {
			ctrl := gomock.NewController(GinkgoT())
			mc := mock_v1.NewMockSecretServiceClient(ctrl)
			s := &secretRefImpl{
				name: "test",
				sc:   mc,
			}

			sv := s.Version("test")

			svi, ok := sv.(*secretVersionRefImpl)

			It("should be of type secretVersionRefImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should share a secret client refernce with its parent secret", func() {
				Expect(svi.sc).To(Equal(s.sc))
			})

			It("should have a back reference to it's parent secret", func() {
				Expect(svi.secret).To(Equal(s))
			})

			It("should have a the requested version name", func() {
				Expect(svi.version).To(Equal("test"))
			})
		})
	})

	Context("Latest", func() {
		When("retrieving a the latest secret version reference", func() {
			ctrl := gomock.NewController(GinkgoT())
			mc := mock_v1.NewMockSecretServiceClient(ctrl)
			s := &secretRefImpl{
				name: "test",
				sc:   mc,
			}

			sv := s.Latest()

			svi, ok := sv.(*secretVersionRefImpl)

			It("should be of type secretVersionRefImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should share a secret client refernce with its parent secret", func() {
				Expect(svi.sc).To(Equal(s.sc))
			})

			It("should have a back reference to it's parent secret", func() {
				Expect(svi.secret).To(Equal(s))
			})

			It("should have 'latest' as it's version name", func() {
				Expect(svi.version).To(Equal("latest"))
			})
		})
	})

	Context("Put", func() {
		When("the RPC server returns successfully", func() {
			ctrl := gomock.NewController(GinkgoT())
			mc := mock_v1.NewMockSecretServiceClient(ctrl)
			s := &secretRefImpl{
				name: "test",
				sc:   mc,
			}

			// NOTE: Using a single It here to correctly assert the number
			// of service calls we are making
			It("should return a reference to the created secret version", func() {
				defer ctrl.Finish()
				By("calling the RPC server with the correct request")
				mc.EXPECT().Put(gomock.Any(), &v1.SecretPutRequest{
					Secret: &v1.Secret{
						Name: "test",
					},
					Value: []byte("ssssshhhh... it's a secret"),
				}).Return(&v1.SecretPutResponse{
					SecretVersion: &v1.SecretVersion{
						Secret: &v1.Secret{
							Name: "test",
						},
						Version: "1",
					},
				}, nil).Times(1)

				// Call the service
				sv, err := s.Put([]byte("ssssshhhh... it's a secret"))

				By("not returning an error")
				Expect(err).ToNot(HaveOccurred())

				svi, ok := sv.(*secretVersionRefImpl)

				By("returning a secretVersionRefImpl")
				Expect(ok).To(BeTrue())

				By("returning the correct secret version")
				Expect(svi.version).To(Equal("1"))

				By("attaching a reference to the parent secret")
				Expect(svi.secret).To(Equal(s))
			})
		})

		When("the RPC server returns an error", func() {
			ctrl := gomock.NewController(GinkgoT())
			mc := mock_v1.NewMockSecretServiceClient(ctrl)
			s := &secretRefImpl{
				name: "test",
				sc:   mc,
			}

			It("should pass through the error", func() {
				defer ctrl.Finish()
				By("calling the RPC server")
				mc.EXPECT().Put(
					gomock.Any(),
					gomock.Any(),
				).Return(nil, fmt.Errorf("mock-error")).Times(1)

				// Call the service
				sv, err := s.Put([]byte("ssssshhhh... it's a secret"))

				By("returning the error")
				Expect(err).To(HaveOccurred())

				By("returning a nil secret version reference")
				Expect(sv).To(BeNil())
			})
		})
	})
})
