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
	JobHandler    = func(*JobContext) (*JobContext, error)
	JobMiddleware = func(*JobContext, JobHandler) (*JobContext, error)
)

func JobDummy(ctx *JobContext) (*JobContext, error) {
	return ctx, nil
}

type chainedJobMiddleware struct {
	fun      JobMiddleware
	nextFunc JobHandler
}

// automatically finalize chain with dummy function
func (c *chainedJobMiddleware) invoke(ctx *JobContext) (*JobContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = JobDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type jobMiddlewareChain struct {
	chain []*chainedJobMiddleware
}

func (h *jobMiddlewareChain) invoke(ctx *JobContext, next JobHandler) (*JobContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeJobMiddleware - Composes an array of middleware into a single middleware
func ComposeJobMiddleware(funcs ...JobMiddleware) JobMiddleware {
	mwareChain := &jobMiddlewareChain{
		chain: make([]*chainedJobMiddleware, len(funcs)),
	}

	var nextFunc JobHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedJobMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
