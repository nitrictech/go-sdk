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

package common

import (
	"encoding/json"

	"github.com/nitrictech/go-sdk/faas"
)

// Json - Middleware parsing http context data as map[string]interface{}
func Json(key string) faas.HttpMiddleware {
	return func(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
		js := make(map[string]interface{})

		if err := json.Unmarshal(ctx.Request.Data(), &js); err != nil {
			ctx.Response.Body = []byte("Bad Request: Expected JSON body")
			ctx.Response.Status = 400

			return ctx, nil
		}

		// decode into the copy of the struct
		ctx.Extras[key] = js

		return next(ctx)
	}
}
