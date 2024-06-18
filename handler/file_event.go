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
	FileEventHandler    = func(*FileEventContext) (*FileEventContext, error)
	FileEventMiddleware = func(*FileEventContext, FileEventHandler) (*FileEventContext, error)
)

func fileEventDummy(ctx *FileEventContext) (*FileEventContext, error) {
	return ctx, nil
}

type chainedFileEventMiddleware struct {
	fun      FileEventMiddleware
	nextFunc FileEventHandler
}

// automatically finalize chain with dummy function
func (c *chainedFileEventMiddleware) invoke(ctx *FileEventContext) (*FileEventContext, error) {
	if c.nextFunc == nil {
		c.nextFunc = fileEventDummy
	}

	return c.fun(ctx, c.nextFunc)
}

type fileEventMiddlewareChain struct {
	chain []*chainedFileEventMiddleware
}

func (h *fileEventMiddlewareChain) invoke(ctx *FileEventContext, next FileEventHandler) (*FileEventContext, error) {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeEventMiddleware - Composes an array of middleware into a single middleware
func ComposeFileEventMiddleware(funcs ...FileEventMiddleware) FileEventMiddleware {
	mwareChain := &fileEventMiddlewareChain{
		chain: make([]*chainedFileEventMiddleware, len(funcs)),
	}

	var nextFunc FileEventHandler = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		cm := &chainedFileEventMiddleware{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
