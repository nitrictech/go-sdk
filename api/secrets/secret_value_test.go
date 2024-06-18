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

var _ = Describe("secretValueImpl", func() {
	var (
		ctrl        *gomock.Controller
		mockSC      *mock_v1.MockSecretManagerClient
		secretName  string
		versionName string
		sv          SecretVersionRef
		secretValue SecretValue
		value       []byte
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSC = mock_v1.NewMockSecretManagerClient(ctrl)
		secretName = "test-secret"
		versionName = "test-version"

		sv = &secretVersionRefImpl{
			secretClient: mockSC,
			secret: &secretRefImpl{
				name:         secretName,
				secretClient: mockSC,
			},
			version: versionName,
		}

		value = []byte("ssssshhhh... it's a secret")
		secretValue = &secretValueImpl{
			version: sv,
			val:     value,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Version", func() {
		It("should return the correct secret version reference", func() {
			Expect(secretValue.Version()).To(Equal(sv))
		})
	})

	Describe("AsBytes", func() {
		It("should return the correct secret value as bytes", func() {
			Expect(secretValue.AsBytes()).To(Equal(value))
		})
	})

	Describe("AsString", func() {
		It("should return the correct secret value as string", func() {
			Expect(secretValue.AsString()).To(Equal(string(value)))
		})
	})
})
