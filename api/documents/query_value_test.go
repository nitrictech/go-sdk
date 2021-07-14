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

var _ = Describe("QueryValue", func() {
	Context("NumberValue", func() {
		When("creating a new NumberValue", func() {
			nv := NumberValue(999)

			It("should have the provided number_value", func() {
				Expect(*nv.number_value).To(Equal(999))
			})

			It("should have nil double_value", func() {
				Expect(nv.double_value).To(BeNil())
			})

			It("should have nil string_value", func() {
				Expect(nv.string_value).To(BeNil())
			})

			It("should have nil bool_value", func() {
				Expect(nv.bool_value).To(BeNil())
			})

			When("translating toWire", func() {
				ev, err := nv.toWire()

				It("should not return an error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should return a proto ExpressionValue", func() {
					Expect(ev.GetIntValue()).To(Equal(int64(999)))
				})
			})
		})
	})

	Context("StringValue", func() {
		When("creating a new StringValue", func() {
			nv := StringValue("test")

			It("should have the provided number_value", func() {
				Expect(*nv.string_value).To(Equal("test"))
			})

			It("should have nil double_value", func() {
				Expect(nv.double_value).To(BeNil())
			})

			It("should have nil number_value", func() {
				Expect(nv.number_value).To(BeNil())
			})

			It("should have nil bool_value", func() {
				Expect(nv.bool_value).To(BeNil())
			})

			When("translating toWire", func() {
				ev, err := nv.toWire()

				It("should not return an error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should return a proto ExpressionValue", func() {
					Expect(ev.GetStringValue()).To(Equal("test"))
				})
			})
		})
	})

	Context("DoubleValue", func() {
		When("creating a new DoubleValue", func() {
			nv := DoubleValue(99.9)

			It("should have the provided number_value", func() {
				Expect(*nv.double_value).To(Equal(99.9))
			})

			It("should have nil string_value", func() {
				Expect(nv.string_value).To(BeNil())
			})

			It("should have nil number_value", func() {
				Expect(nv.number_value).To(BeNil())
			})

			It("should have nil bool_value", func() {
				Expect(nv.bool_value).To(BeNil())
			})

			When("translating toWire", func() {
				ev, err := nv.toWire()

				It("should not return an error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should return a proto ExpressionValue", func() {
					Expect(ev.GetDoubleValue()).To(Equal(float64(99.9)))
				})
			})
		})
	})

	Context("BoolValue", func() {
		When("creating a new BoolValue", func() {
			nv := BoolValue(true)

			It("should have the provided number_value", func() {
				Expect(*nv.bool_value).To(Equal(true))
			})

			It("should have nil string_value", func() {
				Expect(nv.string_value).To(BeNil())
			})

			It("should have nil number_value", func() {
				Expect(nv.number_value).To(BeNil())
			})

			It("should have nil double_value", func() {
				Expect(nv.double_value).To(BeNil())
			})

			When("translating toWire", func() {
				ev, err := nv.toWire()

				It("should not return an error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should return a proto ExpressionValue", func() {
					Expect(ev.GetBoolValue()).To(Equal(true))
				})
			})
		})
	})

	Context("InvalidValue", func() {
		When("Provided an invalid value", func() {
			v := &value{}
			When("translating toWire", func() {
				_, err := v.toWire()

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid query value"))
				})
			})
		})
	})
})
