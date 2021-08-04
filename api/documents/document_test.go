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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Document", func() {

	Context("Content", func() {
		When("given a document with content", func() {
			md := &documentImpl{
				content: map[string]interface{}{
					"test": "test",
				},
			}

			It("should return the documents content", func() {
				Expect(md.Content()).To(Equal(map[string]interface{}{
					"test": "test",
				}))
			})
		})

		When("given a document with nil content", func() {
			It("should return nil content", func() {
				md := &documentImpl{}
				Expect(md.Content()).To(BeNil())
			})
		})
	})

	Context("Decode", func() {
		When("Decoding to a compatible struct", func() {
			type Test struct {
				Test string
			}

			md := &documentImpl{
				content: map[string]interface{}{
					"test": "test",
				},
			}

			var test *Test

			err := md.Decode(&test)

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should populate the provided struct", func() {
				Expect(test).ToNot(BeNil())
				Expect(test.Test).To(Equal("test"))
			})
		})

		When("Decoding an incompatible struct", func() {
			type Test struct {
				Blah string
			}

			md := &documentImpl{
				content: map[string]interface{}{
					"test": "test",
				},
			}

			var test *Test

			err := md.Decode(test)

			It("should not populate test", func() {
				Expect(test).To(BeNil())
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Internal: Document.Decode: \n result must be addressable (a pointer)"))
			})
		})
	})
})
