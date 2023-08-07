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

	"github.com/nitrictech/go-sdk/faas"
)

var _ = Describe("schedule", func() {
	Context("New", func() {
		m := &manager{
			workers:  map[string]Starter{},
			builders: map[string]faas.HandlerBuilder{},
		}
		When("valid args", func() {
			err := m.newSchedule("regular").Every("4 minutes", func(ec *faas.EventContext, eh faas.EventHandler) (*faas.EventContext, error) {
				return eh(ec)
			})

			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
				b := m.builders["regular"]
				Expect(b).ToNot(BeNil())
			})
		})
		When("invalid schedule", func() {
			err := m.newSchedule("invalid").Every("four minutes", func(ec *faas.EventContext, eh faas.EventHandler) (*faas.EventContext, error) {
				return eh(ec)
			})

			It("should return an error", func() {
				Expect(err).To(MatchError("invalid rate expression four minutes; strconv.Atoi: parsing \"four\": invalid syntax"))
				b := m.builders["invalid"]
				Expect(b).To(BeNil())
			})
		})
	})
	Context("rateSplit", func() {
		When("hours", func() {
			It("2 hours", func() {
				r, f, err := rateSplit("2 hours")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(2))
				Expect(f).To(Equal(faas.Frequency("hours")))
			})
			It("1 hours", func() {
				r, f, err := rateSplit("1 hours")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(1))
				Expect(f).To(Equal(faas.Frequency("hours")))
			})
			It("hour", func() {
				r, f, err := rateSplit("hour")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(1))
				Expect(f).To(Equal(faas.Frequency("hours")))
			})
		})
		When("days", func() {
			It("day", func() {
				r, f, err := rateSplit("day")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(1))
				Expect(f).To(Equal(faas.Frequency("days")))
			})
			It("1 day", func() {
				r, f, err := rateSplit("1 days")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(1))
				Expect(f).To(Equal(faas.Frequency("days")))
			})
			It("89 day", func() {
				r, f, err := rateSplit("89 days")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(89))
				Expect(f).To(Equal(faas.Frequency("days")))
			})
		})
		When("minutes", func() {
			It("minute", func() {
				r, f, err := rateSplit("minute")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(1))
				Expect(f).To(Equal(faas.Frequency("minutes")))
			})
			It("1 minutes", func() {
				r, f, err := rateSplit("1 minutes")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(1))
				Expect(f).To(Equal(faas.Frequency("minutes")))
			})
			It("89 minutes", func() {
				r, f, err := rateSplit("89 minutes")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(r).To(Equal(89))
				Expect(f).To(Equal(faas.Frequency("minutes")))
			})
		})
	})
})
