package faas

import (
	"os"

	"github.com/valyala/fasthttp"
)

// NitricFunction - a function built using Nitric, to be executed
type NitricFunction func(*NitricRequest) *NitricResponse

// Start - Starts accepting requests for the provided NitricFunction
//
// This should be the only method called in the 'main' method of your entrypoint package
func Start(f NitricFunction) error {
	var childAddress = "127.0.0.1:8080"
	if env, ok := os.LookupEnv("CHILD_ADDRESS"); ok {
		childAddress = env
	}

	return fasthttp.ListenAndServe(childAddress, func(ctx *fasthttp.RequestCtx) {
		nr := fromRequestContext(ctx)

		response := f(nr)

		// Write the reponse
		response.writeHTTPResponse(ctx)
	})
}
