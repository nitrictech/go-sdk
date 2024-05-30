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

var _ = Describe("secretRefImpl", func() {
	var (
		ctrl     *gomock.Controller
		mockSC  	 *mock_v1.MockSecretManagerClient
		secrets   *secretsImpl
		secretsName 	 string
		sr 		 SecretRef
		ctx			context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSC = mock_v1.NewMockSecretManagerClient(ctrl)
		secrets = &secretsImpl{
			secretClient: mockSC,
		}
		secretsName = "test-secrets"
		sr = secrets.Secret(secretsName)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Having a valid secretRef", func() {
		Describe("Name", func() {
			It("should return the correct secrets name", func(){
				Expect(sr.Name()).To(Equal(secretsName))
			})
		})

		Describe("Version", func() {
			versionName := "test-version"
			_sr := &secretRefImpl{
				name: secretsName,
				secretClient: mockSC,
			}
			sv := _sr.Version(versionName)
			svi, ok := sv.(*secretVersionRefImpl)
		
			It("should be of type secretVersionRefImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should share a secret client references with its parent secret", func() {
				Expect(svi.secretClient).To(Equal(_sr.secretClient))
			})

			It("should have a back reference to it's parent secret", func() {
				Expect(svi.secret).To(Equal(_sr))
			})

			It("should have a the requested version name", func() {
				Expect(svi.version).To(Equal(versionName))
			})
		})
		
		Describe("Latest", func() {
			When("retrieving a the latest secret version reference", func() {
				_sr := &secretRefImpl{
					name: secretsName,
					secretClient: mockSC,
				}
				sv := _sr.Latest()
				svi, ok := sv.(*secretVersionRefImpl)

				It("should be of type secretVersionRefImpl", func() {
					Expect(ok).To(BeTrue())
				})
	
				It("should share a secret client references with its parent secret", func() {
					Expect(svi.secretClient).To(Equal(_sr.secretClient))
				})
	
				It("should have a back reference to it's parent secret", func() {
					Expect(svi.secret).To(Equal(_sr))
				})
	
				It("should have 'latest' as it's version name", func() {
					Expect(svi.version).To(Equal("latest"))
				})
			})
		})

		Describe("Put", func() {

			var(
				_sr  *secretRefImpl
				secretValue []byte
			)

			BeforeEach(func ()  {
				_sr = &secretRefImpl{
					name: secretsName,
					secretClient: mockSC,
				}
				secretValue = []byte("ssssshhhh... it's a secret")
			})

			When("the RPC operation is successful", func ()  {
				var (
					secret *v1.Secret
					versionName string
				)

				BeforeEach(func ()  {
					secret = &v1.Secret{
						Name: _sr.Name(),
					}	
					versionName = "1"

					mockSC.EXPECT().Put(gomock.Any(), &v1.SecretPutRequest{
						Secret: secret,
						Value: secretValue,
					}).Return(
						&v1.SecretPutResponse{
							SecretVersion: &v1.SecretVersion{
								Secret: secret,
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

					svi, ok := sv.(*secretVersionRefImpl)

					By("returning a secretVersionRefImpl")
					Expect(ok).To(BeTrue())

					By("returning the correct secret version")
					Expect(svi.version).To(Equal(versionName))

					By("attaching a reference to the parent secret")
					Expect(svi.secret).To(Equal(_sr))
				})

			})

			When("the RPC operation fails", func ()  {
				var (
					errorMsg string
				)

				BeforeEach(func ()  {
					errorMsg = "Internal Error"
					
					mockSC.EXPECT().Put(gomock.Any(),gomock.Any()).Return(
						nil,
						errors.New(errorMsg),
					).Times(1)
				})

				It("should pass through the error", func() {
					sv, err := _sr.Put(ctx, secretValue)
	
					By("returning the error")
					Expect(err).To(HaveOccurred())
					Expect(strings.Contains(err.Error(), errorMsg)).To(BeTrue())
	
					By("returning a nil secret version reference")
					Expect(sv).To(BeNil())
				})
			})
		})
	})
})