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
	MessageHandler    = func(*MessageContext) (*MessageContext, error)
	MessageMiddleware = func(*MessageContext, MessageHandler) (*MessageContext, error)
)

func messageDummy(ctx *MessageContext) (*MessageContext, error) {
	return ctx, nil
}

type chainedMessageMiddleware struct {
	fun      MessageMiddleware
	nextFunc MessageHandler
}

// automatically finalize chain with dummy function
func (c *chainedMessageMiddleware) invoke(ctx *MessageContext) (*MessageContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = messageDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type messageMiddlewareChain struct {
	chain []*chainedMessageMiddleware
}

func (h *messageMiddlewareChain) invoke(ctx *MessageContext, next MessageHandler) (*MessageContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeMessageMiddleware - Composes an array of middleware into a single middleware
func ComposeMessageMiddleware(funcs ...MessageMiddleware) MessageMiddleware {
	mwareChain := &messageMiddlewareChain{
		chain: make([]*chainedMessageMiddleware, len(funcs)),
	}

	var nextFunc MessageHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedMessageMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
