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

package main

import (
	"fmt"

	"github.com/nitrictech/go-sdk/api/documents"
	"github.com/nitrictech/go-sdk/examples/faas/common"
	"github.com/nitrictech/go-sdk/faas"
)

// [START snippet]
func handler(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
	params, ok := ctx.Extras["params"].(map[string]string)

	if !ok || params == nil {
		return nil, fmt.Errorf("error retrieving path params")
	}

	id := params["id"]

	dc, err := documents.New()
	if err != nil {
		return nil, err
	}

	err = dc.Collection("examples").Doc(id).Delete()
	if err != nil {
		ctx.Response.Body = []byte("Error deleting document")
		ctx.Response.Status = 500
	} else {
		ctx.Response.Body = []byte("Successfully deleted document")
		ctx.Response.Status = 200
	}

	return next(ctx)
}

func main() {
	err := faas.New().Http(
		// Retrieve path parameters if available
		common.PathParser("/examples/:id"),
		// Actual Handler
		handler,
	).Start()

	if err != nil {
		fmt.Println(err)
	}
}

// [END snippet]
