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

type TriggerHandler = func(TriggerContext) (TriggerContext, error)
type TriggerMiddleware = func(TriggerContext, TriggerHandler) (TriggerContext, error)

func triggerDummy(ctx TriggerContext) (TriggerContext, error) {
	return ctx, nil
}

type chainedTriggerMiddleware struct {
	fun      TriggerMiddleware
	nextFunc TriggerHandler
}

// automatically finalize chain with dummy function
func (c *chainedTriggerMiddleware) invoke(ctx TriggerContext) (TriggerContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = triggerDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type triggerMiddlewareChain struct {
	chain []*chainedTriggerMiddleware
}

func (h *triggerMiddlewareChain) invoke(ctx TriggerContext, next TriggerHandler) (TriggerContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// CreateTriggerMiddleware - Chains Trigger middleware functions together to single handler
func ComposeTriggerMiddleware(funcs ...TriggerMiddleware) TriggerMiddleware {
	mwareChain := &triggerMiddlewareChain{
		chain: make([]*chainedTriggerMiddleware, len(funcs)),
	}

	var nextFunc TriggerHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedTriggerMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
