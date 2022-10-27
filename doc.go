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

// Package go-sdk is the Go SDK for the Nitric framework.
//
// Introduction
//
// The Go SDK supports the use of the Nitric framework with Go 1.17+. For more information, check out the main Nitric repo.
//
// Nitric SDKs provide an infrastructure-as-code style that lets you define resources in code. You can also write the functions that support the logic behind APIs, subscribers and schedules.
// You can request the type of access you need to resources such as publishing for topics, without dealing directly with IAM or policy documents.
//
// A good starting point is to check out [resources.NewApi] function.
//
//   exampleApi, err := resources.NewApi("example")
//   if err != nil {
//   	fmt.Println(err)
//   	os.Exit(1)
//   }
//
//   exampleApi.Get("/hello/:name", func(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
//   	params := ctx.Request.PathParams()
//
//   	if params == nil || len(params["name"]) == 0 {
//   		ctx.Response.Body = []byte("error retrieving path params")
//   		ctx.Response.Status = http.StatusBadRequest
//   	} else {
//   		ctx.Response.Body = []byte("Hello " + params["name"])
//   		ctx.Response.Status = http.StatusOK
//   	}
//
//   	return next(ctx)
//   })
//
//   fmt.Println("running example API")
//   if err := resources.Run(); err != nil {
//   	fmt.Println(err)
//   	os.Exit(1)
//   }
//
package main

import _ "github.com/nitrictech/go-sdk/resources"
