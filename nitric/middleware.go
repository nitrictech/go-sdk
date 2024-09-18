package nitric

import "fmt"

type Handler[T any] func(context *T) (*T, error)
type Middleware[T any] func(context *T, next Handler[T]) (*T, error)

type chainedMiddleware[T any] struct {
	fun      Middleware[T]
	nextFunc Handler[T]
}

func (c *chainedMiddleware[T]) invoke(ctx *T) (*T, error) {
	// Chains are left open-ended so middleware can continue to be linked
	// If the chain is incomplete, set a chained dummy handler for safety
	if c.nextFunc == nil {
		c.nextFunc = func(ctx *T) (*T, error) {
			return ctx, nil
		}
	}

	return c.fun(ctx, c.nextFunc)
}

type middlewareChain[T any] struct {
	chain []*chainedMiddleware[T]
}

func (h *middlewareChain[T]) invoke(ctx *T, next Handler[T]) (*T, error) {
	if len(h.chain) == 0 {
		return nil, fmt.Errorf("there are no middleware in this chain")
	}
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

func Compose[T any](funcs ...Middleware[T]) Middleware[T] {
	mwareChain := &middlewareChain[T]{
		chain: make([]*chainedMiddleware[T], len(funcs)),
	}

	var nextFunc Handler[T] = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		if funcs[i] == nil {
			fmt.Println("this func is empty")
		}

		cm := &chainedMiddleware[T]{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain[i] = cm
	}

	return mwareChain.invoke
}
