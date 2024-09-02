package job

type (
	Middleware = func(*Context, func(*Context) error) error
	Next       = func(*Context) error
)

func NoOpHandler(ctx *Context) error {
	return nil
}

type chainedJobMiddleware struct {
	fun      Middleware
	nextFunc Next
}

// automatically finalize chain with dummy function
func (c *chainedJobMiddleware) invoke(ctx *Context) error {
	if c.nextFunc == nil {
		c.nextFunc = NoOpHandler
	}

	return c.fun(ctx, c.nextFunc)
}

type jobMiddlewareChain struct {
	chain []*chainedJobMiddleware
}

func (h *jobMiddlewareChain) invoke(ctx *Context, next Next) error {
	// Complete the chain
	h.chain[len(h.chain)-1].nextFunc = next

	return h.chain[0].invoke(ctx)
}

// ComposeJobMiddleware - Composes an array of middleware into a single middleware
func ComposeJobMiddleware(funcs ...Middleware) Middleware {
	mwareChain := &jobMiddlewareChain{
		chain: make([]*chainedJobMiddleware, len(funcs)),
	}

	var nextFunc Next = nil
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
