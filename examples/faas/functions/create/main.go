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

	"github.com/google/uuid"
	"github.com/nitrictech/go-sdk/api/documents"
	"github.com/nitrictech/go-sdk/examples/faas/common"
	"github.com/nitrictech/go-sdk/faas"
)

// [START snippet]
const exampleKey = "example"

func handler(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
	id := uuid.New().String()
	example, ok := ctx.Extras[exampleKey].(map[string]interface{})

	if !ok || example == nil {
		return nil, fmt.Errorf("unable to retrieve decoded example")
	}

	dc, err := documents.New()

	if err != nil {
		return nil, err
	}

	if err := dc.Collection("examples").Doc(id).Set(example); err != nil {
		return nil, err
	}

	ctx.Response.Status = 200
	ctx.Response.Body = []byte(fmt.Sprintf("Created example with ID: %s", id))

	return next(ctx)
}

func main() {
	err := faas.New().Http(
		// Decoding middleware
		common.Json(exampleKey),
		// Actual Handler
		handler,
	).Start()

	if err != nil {
		fmt.Println(err)
	}
}

// [END snippet]
