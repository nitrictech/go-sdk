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

var _ = Describe("secretVersionRefImpl", func() {
	var (
		ctrl        *gomock.Controller
		mockSC      *mock_v1.MockSecretManagerClient
		secretName  string
		versionName string
		sv          SecretVersionRef
		sr 		SecretRef
		ctx         context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSC = mock_v1.NewMockSecretManagerClient(ctrl)
		secretName = "test-secret"
		versionName = "test-version"
		
		sr = &secretRefImpl{
			name:         secretName,
			secretClient: mockSC,
		}

		sv = &secretVersionRefImpl{
			secretClient: mockSC,
			secret: sr,
			version: versionName,
		}
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Access", func() {
		var secretValue []byte

		BeforeEach(func() {
			secretValue = []byte("super-secret-value")
		})

		When("the RPC operation is successful", func() {
			BeforeEach(func() {
				mockSC.EXPECT().Access(gomock.Any(), &v1.SecretAccessRequest{
					SecretVersion: &v1.SecretVersion{
						Secret: &v1.Secret{
							Name: secretName,
						},
						Version: versionName,
					},
				}).Return(
					&v1.SecretAccessResponse{
						SecretVersion: &v1.SecretVersion{
							Secret: &v1.Secret{
								Name: secretName,
							},
							Version: versionName,
						},
						Value: secretValue,
					}, nil,
				)
			})

			It("should return the secret value", func() {
				svValue, err := sv.Access(ctx)

				By("not returning an error")
				Expect(err).ToNot(HaveOccurred())

				By("returning a SecretValue")
				Expect(svValue).ToNot(BeNil())

				By("returning the correct secret value")
				Expect(svValue.AsBytes()).To(Equal(secretValue))

				By("returning the correct secret version")
				Expect(svValue.Version().Version()).To(Equal(versionName))
			})
		})

		When("the RPC operation fails", func() {
			var errorMsg string

			BeforeEach(func() {
				errorMsg = "Internal Error"
				mockSC.EXPECT().Access(gomock.Any(), gomock.Any()).Return(
					nil,
					errors.New(errorMsg),
				).Times(1)
			})

			It("should return an error", func() {
				svValue, err := sv.Access(ctx)

				By("returning the error")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())

				By("returning a nil SecretValue")
				Expect(svValue).To(BeNil())
			})
		})
	})

	Context("Version", func() {
		When("retrieving the version of a secretVersionRefImpl", func() {
			It("should return it's internal version field", func() {
				Expect(sv.Version()).To(Equal(versionName))
			})
		})
	})

	
	Context("Secret", func() {
		When("retrieving the parent secret of a secretVersionRefImpl", func() {
			It("should return it's internal secret field", func() {
				Expect(sv.Secret()).To(Equal(sr))
			})
		})
	})
})