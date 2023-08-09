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

type (
	BucketNotificationHandler    = func(*BucketNotificationContext) (*BucketNotificationContext, error)
	BucketNotificationMiddleware = func(*BucketNotificationContext, BucketNotificationHandler) (*BucketNotificationContext, error)
)

func bucketNotificationDummy(ctx *BucketNotificationContext) (*BucketNotificationContext, error) {
	return ctx, nil
}

type chainedBucketNotificationMiddleware struct {
	fun      BucketNotificationMiddleware
	nextFunc BucketNotificationHandler
}

// automatically finalize chain with dummy function
func (c *chainedBucketNotificationMiddleware) invoke(ctx *BucketNotificationContext) (*BucketNotificationContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = bucketNotificationDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type bucketNotificationMiddlewareChain struct {
	chain []*chainedBucketNotificationMiddleware
}

func (h *bucketNotificationMiddlewareChain) invoke(ctx *BucketNotificationContext, next BucketNotificationHandler) (*BucketNotificationContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeEventMiddleware - Composes an array of middleware into a single middleware
func ComposeBucketNotificationMiddleware(funcs ...BucketNotificationMiddleware) BucketNotificationMiddleware {
	mwareChain := &bucketNotificationMiddlewareChain{
		chain: make([]*chainedBucketNotificationMiddleware, len(funcs)),
	}

	var nextFunc BucketNotificationHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedBucketNotificationMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
