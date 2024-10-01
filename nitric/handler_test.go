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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("handler", func() {
	// func()
	// func() error
	// func(*T)
	// func(*T) error
	// func(*T, Handler[T]) error
	// Handler[T]
	Context("interfaceToHandler", func() {
		When("interface{} is func()", func() {
			It("should return a valid handler", func() {
				handler, err := interfaceToHandler[any](func() {})

				Expect(err).To(BeNil())
				Expect(handler).ToNot(BeNil())
			})
		})

		When("interface{} is func() error", func() {
			It("should return a valid handler", func() {
				handler, err := interfaceToHandler[any](func() error { return nil })

				Expect(err).To(BeNil())
				Expect(handler).ToNot(BeNil())
			})
		})

		When("interface{} is func(*T)", func() {
			It("should return a valid handler", func() {
				handler, err := interfaceToHandler[string](func(*string) {})

				Expect(err).To(BeNil())
				Expect(handler).ToNot(BeNil())
			})
		})

		When("interface{} is func(*T) error", func() {
			It("should return a valid handler", func() {
				handler, err := interfaceToHandler[string](func(*string) error { return nil })

				Expect(err).To(BeNil())
				Expect(handler).ToNot(BeNil())
			})
		})

		When("interface{} is not a valid type", func() {
			It("should return an error", func() {
				handler, err := interfaceToHandler[string](func() (error, error) { return nil, nil })

				Expect(err).ToNot(BeNil())
				Expect(handler).To(BeNil())
			})
		})
	})
})
