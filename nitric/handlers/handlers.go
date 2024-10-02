// Copyright 2023 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import "fmt"

type (
	Handler[T any]    func(context *T) error
	Middleware[T any] func(Handler[T]) Handler[T]
)

// fromInterface - Converts a function to a Handler
// Valid function types are:
// func()
// func() error
// func(*T)
// func(*T) error
// Handler[T]
// If the function is not a valid type, an error is returned
func HandlerFromInterface[T any](handler interface{}) (Handler[T], error) {
	var typedHandler Handler[T]
	switch handlerType := handler.(type) {
	case func():
		typedHandler = func(ctx *T) error {
			handlerType()
			return nil
		}
	case func() error:
		typedHandler = func(ctx *T) error {
			return handlerType()
		}
	case func(*T):
		typedHandler = func(ctx *T) error {
			handlerType(ctx)
			return nil
		}
	case func(*T) error:
		typedHandler = Handler[T](handlerType)
	case Handler[T]:
		typedHandler = handlerType

	default:
		return nil, fmt.Errorf("invalid handler type: %T", handler)
	}

	return typedHandler, nil
}
