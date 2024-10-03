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
	"context"
	"errors"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
)

var _ = Describe("Secret", func() {
	var (
		ctrl   *gomock.Controller
		mockSC *mock_v1.MockSecretManagerClient
		secret *SecretClient
		ctx    context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSC = mock_v1.NewMockSecretManagerClient(ctrl)
		secret = &SecretClient{
			name:         "test-secret",
			secretClient: mockSC,
		}
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Having a valid secretRef", func() {
		Describe("Name", func() {
			It("should return the correct secrets name", func() {
				Expect(secret.Name()).To(Equal("test-secret"))
			})
		})

		Describe("Put", func() {
			var (
				_sr         *SecretClient
				secretValue []byte
			)

			BeforeEach(func() {
				_sr = &SecretClient{
					name:         "test-secret",
					secretClient: mockSC,
				}
				secretValue = []byte("ssssshhhh... it's a secret")
			})

			When("the RPC operation is successful", func() {
				var (
					secret      *v1.Secret
					versionName string
				)

				BeforeEach(func() {
					secret = &v1.Secret{
						Name: _sr.Name(),
					}
					versionName = "1"

					mockSC.EXPECT().Put(gomock.Any(), &v1.SecretPutRequest{
						Secret: secret,
						Value:  secretValue,
					}).Return(
						&v1.SecretPutResponse{
							SecretVersion: &v1.SecretVersion{
								Secret:  secret,
								Version: versionName,
							},
						},
						nil,
					).Times(1)
				})

				It("should return a reference to the created secret version", func() {
					sv, err := _sr.Put(ctx, secretValue)

					By("not returning an error")
					Expect(err).ToNot(HaveOccurred())

					By("returning a string secret version")
					Expect(sv).To(Equal(versionName))
				})
			})

			When("the RPC operation fails", func() {
				var errorMsg string

				BeforeEach(func() {
					errorMsg = "Internal Error"

					mockSC.EXPECT().Put(gomock.Any(), gomock.Any()).Return(
						nil,
						errors.New(errorMsg),
					).Times(1)
				})

				It("should pass through the error", func() {
					sv, err := _sr.Put(ctx, secretValue)

					By("returning the error")
					Expect(err).To(HaveOccurred())
					Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())

					By("returning a blank secret version reference")
					Expect(sv).To(BeEmpty())
				})
			})
		})
	})
})
