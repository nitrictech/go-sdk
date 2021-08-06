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

var _ = Describe("QueryExpression", func() {
	Context("QueryOp", func() {
		Context("IsValid", func() {
			When("given a valid QueryOp", func() {
				It("should not return an error", func() {
					Expect(queryOp_EQ.IsValid()).ToNot(HaveOccurred())
				})
			})

			When("given an invalid QueryOp", func() {
				var test queryOp = "test"

				err := test.IsValid()

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid Argument: queryOp.IsValid: invalid query operation (test)"))
				})
			})
		})
	})

	Context("toWire", func() {
		When("translating a valid query expression to wire", func() {
			qe := &queryExpression{
				field: "test",
				op:    queryOp_EQ,
				val:   StringValue("test"),
			}

			r, err := qe.toWire()

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should contain the provided name", func() {
				Expect(r.GetOperand()).To(Equal("test"))
			})

			It("should contain the provider operation", func() {
				Expect(r.GetOperator()).To(Equal("=="))
			})

			It("should contain the provided value", func() {
				Expect(r.GetValue().GetStringValue()).To(Equal("test"))
			})
		})

		When("translating an query expression missing its field", func() {
			qe := &queryExpression{
				op:  queryOp_EQ,
				val: StringValue("test"),
			}

			_, err := qe.toWire()

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid Argument: queryExpress.toWire: provide non-blank field name"))
			})
		})

		When("translating a query expression with an invalid op", func() {
			qe := &queryExpression{
				field: "test",
				op:    "blah",
				val:   StringValue("test"),
			}

			_, err := qe.toWire()

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid Argument: queryOp.IsValid: invalid query operation (blah)"))
			})
		})

		When("translating a valid query expression with an invalid value", func() {
			qe := &queryExpression{
				field: "test",
				op:    queryOp_EQ,
				val:   &value{},
			}

			_, err := qe.toWire()

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid Argument: value.toWire: Invalid Query Value"))
			})
		})
	})

	Context("Condition", func() {
		When("creating a new query condition", func() {
			c := Condition("test")

			It("should contain the provided field name", func() {
				Expect(c.field).To(Equal("test"))
			})
		})
	})

	Context("QueryExpressionBuilder", func() {
		c := Condition("test")

		When("creating a new EQ condition", func() {
			qe := c.Eq(StringValue("test"))

			It("should contain the provided field name", func() {
				Expect(qe.field).To(Equal("test"))
			})

			It("should contain the provided value", func() {
				Expect(qe.val).To(Equal(StringValue("test")))
			})

			It("should contain the provided operation", func() {
				Expect(qe.op).To(Equal(queryOp_EQ))
			})
		})

		When("creating a new LT condition", func() {
			qe := c.Lt(NumberValue(999))

			It("should contain the provided field name", func() {
				Expect(qe.field).To(Equal("test"))
			})

			It("should contain the provided value", func() {
				Expect(qe.val).To(Equal(NumberValue(999)))
			})

			It("should contain the provided operation", func() {
				Expect(qe.op).To(Equal(queryOp_LT))
			})
		})

		When("creating a new LE condition", func() {
			qe := c.Le(NumberValue(999))

			It("should contain the provided field name", func() {
				Expect(qe.field).To(Equal("test"))
			})

			It("should contain the provided value", func() {
				Expect(qe.val).To(Equal(NumberValue(999)))
			})

			It("should contain the provided operation", func() {
				Expect(qe.op).To(Equal(queryOp_LE))
			})
		})

		When("creating a new GT condition", func() {
			qe := c.Gt(NumberValue(999))

			It("should contain the provided field name", func() {
				Expect(qe.field).To(Equal("test"))
			})

			It("should contain the provided value", func() {
				Expect(qe.val).To(Equal(NumberValue(999)))
			})

			It("should contain the provided operation", func() {
				Expect(qe.op).To(Equal(queryOp_GT))
			})
		})

		When("creating a new GE condition", func() {
			qe := c.Ge(NumberValue(999))

			It("should contain the provided field name", func() {
				Expect(qe.field).To(Equal("test"))
			})

			It("should contain the provided value", func() {
				Expect(qe.val).To(Equal(NumberValue(999)))
			})

			It("should contain the provided operation", func() {
				Expect(qe.op).To(Equal(queryOp_GE))
			})
		})

		When("creating a new StartsWith condition", func() {
			qe := c.StartsWith(StringValue("test"))

			It("should contain the provided field name", func() {
				Expect(qe.field).To(Equal("test"))
			})

			It("should contain the provided value", func() {
				Expect(qe.val).To(Equal(StringValue("test")))
			})

			It("should contain the provided operation", func() {
				Expect(qe.op).To(Equal(queryOp_StartsWith))
			})
		})
	})
})
