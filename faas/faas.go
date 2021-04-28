// Copyright 2021 Nitric Pty Ltd.
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
