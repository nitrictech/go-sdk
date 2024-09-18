package nitric

import "fmt"

type Handler[T any] func(context *T) (*T, error)
type Middleware[T any] func(context *T, next Handler[T]) (*T, error)

type chainedMiddleware[T any] struct {
	fun      Middleware[T]
	nextFunc Handler[T]
}

func DummyHandler[T any](ctx *T) (*T, error) {
	return ctx, nil
}

func (c *chainedMiddleware[T]) invoke(ctx *T) (*T, error) {
	// Chains are left open-ended so middleware can continue to be linked
	// If the chain is incomplete, set a chained dummy handler for safety
	if c.nextFunc == nil {
		c.nextFunc = DummyHandler[T]
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

// Compose - Takes a collection of middleware and composes it into a single middleware function
func Compose[T any](funcs ...Middleware[T]) Middleware[T] {
	mwareChain := &middlewareChain[T]{
		chain: []*chainedMiddleware[T]{},
	}

	var nextFunc Handler[T] = nil
	for i := len(funcs) - 1; i >= 0; i = i - 1 {
		if funcs[i] == nil {
			continue
		}

		cm := &chainedMiddleware[T]{
			fun:      funcs[i],
			nextFunc: nextFunc,
		}
		nextFunc = cm.invoke
		mwareChain.chain = append(mwareChain.chain, cm)
	}

	return mwareChain.invoke
}
