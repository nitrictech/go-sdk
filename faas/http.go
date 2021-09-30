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

type HttpHandler = func(*HttpContext) (*HttpContext, error)
type HttpMiddleware = func(*HttpContext, HttpHandler) (*HttpContext, error)

func httpDummy(ctx *HttpContext) (*HttpContext, error) {
	return ctx, nil
}

type chainedHttpMiddleware struct {
	fun      HttpMiddleware
	nextFunc HttpHandler
}

// automatically finalize chain with dummy function
func (c *chainedHttpMiddleware) invoke(ctx *HttpContext) (*HttpContext, error) {
	// Chains are left open-ended so middleware can continue to be linked
	// If the chain is incomplete, set a chained dummy handler for safety
	if c.nextFunc == nil {
		c.nextFunc = httpDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type httpMiddlewareChain struct {
	chain []*chainedHttpMiddleware
}

func (h *httpMiddlewareChain) invoke(ctx *HttpContext, next HttpHandler) (*HttpContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// CreateHttpHandler - Chains HTTP middleware functions together to single handler
func CreateHttpHandler(funcs ...HttpMiddleware) HttpMiddleware {
	mwareChain := &httpMiddlewareChain{
		chain: make([]*chainedHttpMiddleware, len(funcs)),
	}

	var nextFunc HttpHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedHttpMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
