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
	"strings"

	"github.com/nitrictech/go-sdk/faas"
)

const paramToken = ":"
const paramsKey = "params"

// PathParser - Middleware for parsing parameters from HTTP path
func PathParser(paramExpression string) faas.HttpMiddleware {
	pathParts := strings.Split(paramExpression, "/")
	parts := make(map[int]string)

	for i, s := range pathParts {
		if strings.HasPrefix(s, paramToken) {
			parts[i] = strings.Replace(s, paramToken, "", -1)
		}
	}

	return func(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
		pathParts := strings.Split(ctx.Request.Path(), "/")

		params := make(map[string]string)

		for i, part := range pathParts {
			if parts[i] != "" {
				params[parts[i]] = part
			}
		}

		// decode into the copy of the struct
		ctx.Extras["params"] = params

		return next(ctx)
	}
}
