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

package context

type (
	IntervalHandler    = func(*IntervalContext) (*IntervalContext, error)
	IntervalMiddleware = func(*IntervalContext, IntervalHandler) (*IntervalContext, error)
)

func intervalDummy(ctx *IntervalContext) (*IntervalContext, error) {
	return ctx, nil
}

type chainedIntervalMiddleware struct {
	fun      IntervalMiddleware
	nextFunc IntervalHandler
}

// automatically finalize chain with dummy function
func (c *chainedIntervalMiddleware) invoke(ctx *IntervalContext) (*IntervalContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = intervalDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type intervalMiddlewareChain struct {
	chain []*chainedIntervalMiddleware
}

func (h *intervalMiddlewareChain) invoke(ctx *IntervalContext, next IntervalHandler) (*IntervalContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeIntervalMiddleware - Composes an array of middleware into a single middleware
func ComposeIntervalMiddleware(funcs ...IntervalMiddleware) IntervalMiddleware {
	mwareChain := &intervalMiddlewareChain{
		chain: make([]*chainedIntervalMiddleware, len(funcs)),
	}

	var nextFunc IntervalHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedIntervalMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
