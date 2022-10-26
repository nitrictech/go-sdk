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

package resources

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/nitrictech/go-sdk/api/queues"
	"github.com/nitrictech/go-sdk/faas"
)

// This shows how to create an API and use the Get method to register a handler.
func ExampleNewApi() {
	exampleApi, err := NewApi("example")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exampleApi.Get("/hello/:name", func(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
		params := ctx.Request.PathParams()

		if params == nil || len(params["name"]) == 0 {
			ctx.Response.Body = []byte("error retrieving path params")
			ctx.Response.Status = http.StatusBadRequest
		} else {
			ctx.Response.Body = []byte("Hello " + params["name"])
			ctx.Response.Status = http.StatusOK
		}

		return next(ctx)
	})

	fmt.Println("running example API")
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// This shows how to create a Topic and subscribe to it using an event handler.
func ExampleNewTopic() {
	exampleTopic, err := NewTopic("example")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exampleTopic.Subscribe(func(ec *faas.EventContext, next faas.EventHandler) (*faas.EventContext, error) {
		fmt.Printf("event received %s %v", ec.Request.Topic(), string(ec.Request.Data()))

		return next(ec)
	})

	fmt.Println("running example")
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// This shows how to create a Schedule and subscribe to it.
// It also shows how to process a queue within this handler.
func ExampleNewSchedule() {
	queue, err := NewQueue("work", QueueReceving)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = NewSchedule("job", "10 minutes", func(ec *faas.EventContext, next faas.EventHandler) (*faas.EventContext, error) {
		fmt.Println("got scheduled event ", string(ec.Request.Data()))
		tasks, err := queue.Receive(10)
		if err != nil {
			fmt.Println(err)
			return nil, err
		} else {
			for _, task := range tasks {
				fmt.Printf("processing task %s", task.Task().ID)

				if err = task.Complete(); err != nil {
					fmt.Println(err)
				}
			}
		}

		return next(ec)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("running example")
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// See ExampleNewSchedule() for processing a Queue.
func ExampleNewQueue() {
	queue, err := NewQueue("work", QueueSending)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exampleApi, err := NewApi("example")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exampleApi.Get("/hello/:name", func(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
		params := ctx.Request.PathParams()

		if params == nil || len(params["name"]) == 0 {
			ctx.Response.Body = []byte("error retrieving path params")
			ctx.Response.Status = http.StatusBadRequest
		} else {
			_, err = queue.Send([]*queues.Task{
				{
					ID:          uuid.NewString(),
					PayloadType: "custom-X",
					Payload:     map[string]interface{}{"fruit": "apples"},
				},
			})
			if err != nil {
				ctx.Response.Body = []byte("error sending event")
				ctx.Response.Status = http.StatusInternalServerError
			} else {
				ctx.Response.Body = []byte("Hello " + params["name"])
				ctx.Response.Status = http.StatusOK
			}
		}

		return next(ctx)
	})

	fmt.Println("running example")
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
