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

package handler

type (
	WebsocketHandler    = func(*WebsocketContext) (*WebsocketContext, error)
	WebsocketMiddleware = func(*WebsocketContext, WebsocketHandler) (*WebsocketContext, error)
)

func websocketDummy(ctx *WebsocketContext) (*WebsocketContext, error) {
	return ctx, nil
}

type chainedWebsocketMiddleware struct {
	fun      WebsocketMiddleware
	nextFunc WebsocketHandler
}

// automatically finalize chain with dummy function
func (c *chainedWebsocketMiddleware) invoke(ctx *WebsocketContext) (*WebsocketContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = websocketDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type websocketMiddlewareChain struct {
	chain []*chainedWebsocketMiddleware
}

func (h *websocketMiddlewareChain) invoke(ctx *WebsocketContext, next WebsocketHandler) (*WebsocketContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeWebsocketMiddleware - Composes an array of middleware into a single middleware
func ComposeWebsocketMiddleware(funcs ...WebsocketMiddleware) WebsocketMiddleware {
	mwareChain := &websocketMiddlewareChain{
		chain: make([]*chainedWebsocketMiddleware, len(funcs)),
	}

	var nextFunc WebsocketHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedWebsocketMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
