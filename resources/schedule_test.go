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

package resources

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/go-sdk/faas"
)

var _ = Describe("schedule", func() {
	Context("New", func() {
		m := &manager{
			blockers: map[string]Starter{},
			builders: map[string]faas.HandlerBuilder{},
		}
		When("valid args", func() {
			err := m.NewSchedule("regular", "4 minutes", func(ec *faas.EventContext, eh faas.EventHandler) (*faas.EventContext, error) {
				return eh(ec)
			})

			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
				b := m.builders["regular"]
				Expect(b).ToNot(BeNil())
			})
		})
		When("invalid schedule", func() {
			err := m.NewSchedule("invalid", "four minutes", func(ec *faas.EventContext, eh faas.EventHandler) (*faas.EventContext, error) {
				return eh(ec)
			})

			It("should return an error", func() {
				Expect(err).To(MatchError("invalid rate expression four minutes; strconv.Atoi: parsing \"four\": invalid syntax"))
				b := m.builders["invalid"]
				Expect(b).To(BeNil())
			})
		})
	})
})
