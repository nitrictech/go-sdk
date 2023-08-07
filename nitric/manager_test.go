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

package nitric

import (
	"errors"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	apierrors "github.com/nitrictech/go-sdk/api/errors"
)

var _ = Describe("manager", func() {
	Context("isEOF", func() {
		When("is nil", func() {
			got := isEOF(nil)
			It("should be false", func() {
				Expect(got).To(Equal(false))
			})
		})

		When("is EOF wrapped in ApiError", func() {
			err := apierrors.NewWithCause(500, "unknown", io.EOF)
			got := isEOF(err)
			It("should be true", func() {
				Expect(got).To(Equal(true))
			})
		})

		When("is an ApiError with nil cause", func() {
			err := apierrors.New(500, "unknown")
			got := isEOF(err)
			It("should be false", func() {
				Expect(got).To(Equal(false))
			})
		})

		When("is EOF string wrapped in ApiError", func() {
			err := apierrors.NewWithCause(500, "unknown", errors.New("EOF"))
			got := isEOF(err)
			It("should be true", func() {
				Expect(got).To(Equal(true))
			})
		})

		When("is unexpectedEOF wrapped in ApiError", func() {
			err := apierrors.NewWithCause(500, "unknown", io.ErrUnexpectedEOF)
			got := isEOF(err)
			It("should be false", func() {
				Expect(got).To(Equal(false))
			})
		})

		When("is native EOF", func() {
			got := isEOF(io.EOF)
			It("should be true", func() {
				Expect(got).To(Equal(true))
			})
		})
	})
})
