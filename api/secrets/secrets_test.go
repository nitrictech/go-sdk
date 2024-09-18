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
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_v1 "github.com/nitrictech/go-sdk/mocks"
)

var _ = Describe("Secrets API", func() {
	var (
		ctrl    *gomock.Controller
		mockSC  *mock_v1.MockSecretManagerClient
		secrets *secretsImpl
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSC = mock_v1.NewMockSecretManagerClient(ctrl)
		secrets = &secretsImpl{
			secretClient: mockSC,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Secret method", func() {
		When("creating a new Secret reference", func() {
			var (
				sr          SecretRef
				secretsName string
				srImpl      *secretRefImpl
				ok          bool
			)

			BeforeEach(func() {
				secretsName = "test-secret"
				sr = secrets.Secret(secretsName)
				srImpl, ok = sr.(*secretRefImpl)
			})

			It("should be an instance of secretsImpl", func() {
				Expect(ok).To(BeTrue())
			})

			It("should have the provided secrets name", func() {
				Expect(sr.Name()).To(Equal(secretsName))
			})

			It("should share the Secret's gRPC client", func() {
				Expect(srImpl.secretClient).To(Equal(mockSC))
			})
		})
	})
})
