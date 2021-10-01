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

type EventHandler = func(*EventContext) (*EventContext, error)
type EventMiddleware = func(*EventContext, EventHandler) (*EventContext, error)

func eventDummy(ctx *EventContext) (*EventContext, error) {
	return ctx, nil
}

type chainedEventMiddleware struct {
	fun      EventMiddleware
	nextFunc EventHandler
}

// automatically finalize chain with dummy function
func (c *chainedEventMiddleware) invoke(ctx *EventContext) (*EventContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = eventDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type eventMiddlewareChain struct {
	chain []*chainedEventMiddleware
}

func (h *eventMiddlewareChain) invoke(ctx *EventContext, next EventHandler) (*EventContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeEventMiddleware - Composes an array of middleware into a single middleware
func ComposeEventMiddleware(funcs ...EventMiddleware) EventMiddleware {
	mwareChain := &eventMiddlewareChain{
		chain: make([]*chainedEventMiddleware, len(funcs)),
	}

	var nextFunc EventHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedEventMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
