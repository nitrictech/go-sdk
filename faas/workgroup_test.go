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

package faas

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorkPool", func() {
	Context("3 max", func() {

		When("no jobs are provided", func() {
			w := NewWorkPool(3)
			w.Wait()
			Expect(w.Err()).ShouldNot(HaveOccurred())
		})

		When("one error job is run", func() {
			w := NewWorkPool(3)
			w.Go(func(a interface{}) error { return errors.New(a.(string)) }, "one")
			w.Wait()
			Expect(w.Err().Error()).To(Equal("one"))
		})

		When("five jobs are run, be 2 error", func() {
			w := NewWorkPool(3)
			w.Go(func(a interface{}) error { return nil }, nil)
			w.Go(func(a interface{}) error { return errors.New(a.(string)) }, "x")
			w.Go(func(a interface{}) error { return nil }, nil)
			w.Go(func(a interface{}) error { return nil }, nil)
			w.Go(func(a interface{}) error {
				time.Sleep(50 * time.Millisecond)
				return errors.New(a.(string))
			}, "y")
			w.Wait()
			Expect(w.Err().Error()).To(Equal("the following errors occured:\nx\ny\n")) //nolint:misspell
		})
	})
})
