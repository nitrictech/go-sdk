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

package faas

// ResponseContext
type ResponseContext struct {
	context interface{}
}

func (c *ResponseContext) IsHttp() bool {
	_, ok := c.context.(*HttpResponseContext)

	return ok
}

func (c *ResponseContext) AsHttp() *HttpResponseContext {
	if ctx, ok := c.context.(*HttpResponseContext); ok {
		return ctx
	}

	return nil
}

func (c *ResponseContext) IsTopic() bool {
	_, ok := c.context.(*TopicResponseContext)

	return ok
}

func (c *ResponseContext) AsTopic() *TopicResponseContext {
	if ctx, ok := c.context.(*TopicResponseContext); ok {
		return ctx
	}

	return nil
}

type HttpResponseContext struct {
	Headers map[string][]string
	Status  int
}

type TopicResponseContext struct {
	Success bool
}
