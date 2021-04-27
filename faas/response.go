package faas

import (
	"github.com/valyala/fasthttp"
)

// NitricResponse - represents the results of calling a function.
type NitricResponse struct {
	Headers map[string]string
	Status  int
	Body    []byte
}

// writeHttpResponse - writes a HTTP response from a NitricResponse
func (n *NitricResponse) writeHTTPResponse(ctx *fasthttp.RequestCtx) {
	for k, v := range n.Headers {
		ctx.Response.Header.Add(k, v)
	}

	ctx.SetStatusCode(n.Status)
	ctx.SetBody(n.Body)
}
