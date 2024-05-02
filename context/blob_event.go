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
	BlobEventHandler    = func(*BlobEventContext) (*BlobEventContext, error)
	BlobEventMiddleware = func(*BlobEventContext, BlobEventHandler) (*BlobEventContext, error)
)

func blobEventDummy(ctx *BlobEventContext) (*BlobEventContext, error) {
	return ctx, nil
}

type chainedBlobEventMiddleware struct {
	fun      BlobEventMiddleware
	nextFunc BlobEventHandler
}

// automatically finalize chain with dummy function
func (c *chainedBlobEventMiddleware) invoke(ctx *BlobEventContext) (*BlobEventContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = blobEventDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type blobEventMiddlewareChain struct {
	chain []*chainedBlobEventMiddleware
}

func (h *blobEventMiddlewareChain) invoke(ctx *BlobEventContext, next BlobEventHandler) (*BlobEventContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeEventMiddleware - Composes an array of middleware into a single middleware
func ComposeBlobEventMiddleware(funcs ...BlobEventMiddleware) BlobEventMiddleware {
	mwareChain := &blobEventMiddlewareChain{
		chain: make([]*chainedBlobEventMiddleware, len(funcs)),
	}

	var nextFunc BlobEventHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedBlobEventMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
