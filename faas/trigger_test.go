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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Trigger", func() {
	Context("ComposeTriggerMiddleware", func() {

		When("Creating a function with a single handler", func() {
			hndlr := ComposeTriggerMiddleware(func(ctx TriggerContext, next TriggerHandler) (TriggerContext, error) {
				ctx.Http().Response.Status = 201

				return next(ctx)
			})

			It("Should call the provided function", func() {
				ctx, err := hndlr(&triggerContextImpl{
					http: &HttpContext{
						Response: &HttpResponse{
							Status: 200,
						},
					},
				}, nil)

				Expect(err).ToNot(HaveOccurred())
				Expect(ctx.Http().Response.Status).To(BeEquivalentTo(201))
			})
		})

		When("Creating a function from multiple handlers", func() {
			callOrder := make([]string, 0)

			hndlr := ComposeTriggerMiddleware(
				func(ctx TriggerContext, next TriggerHandler) (TriggerContext, error) {
					callOrder = append(callOrder, "1")
					return next(ctx)
				},
				func(ctx TriggerContext, next TriggerHandler) (TriggerContext, error) {
					callOrder = append(callOrder, "2")
					return ctx, nil
				},
			)

			It("Should call the functions in the provided order", func() {
				hndlr(&triggerContextImpl{}, nil)

				Expect(callOrder).To(BeEquivalentTo([]string{"1", "2"}))
			})
		})

		When("Creating a function from multiple nested middlewares", func() {
			callOrder := make([]string, 0)

			hndlr := ComposeTriggerMiddleware(ComposeTriggerMiddleware(
				func(ctx TriggerContext, next TriggerHandler) (TriggerContext, error) {
					callOrder = append(callOrder, "1")
					return next(ctx)
				},
				func(ctx TriggerContext, next TriggerHandler) (TriggerContext, error) {
					callOrder = append(callOrder, "2")
					return next(ctx)
				},
			), ComposeTriggerMiddleware(
				func(ctx TriggerContext, next TriggerHandler) (TriggerContext, error) {
					callOrder = append(callOrder, "3")
					return ctx, nil
				},
			))

			It("Should call the functions in the provided order", func() {
				hndlr(&triggerContextImpl{
					http: &HttpContext{},
				}, nil)

				Expect(callOrder).To(BeEquivalentTo([]string{"1", "2", "3"}))
			})
		})
	})
})
