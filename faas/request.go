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

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// NitricRequest - represents a request to trigger a function, with payload and context required to execute that function.
type NitricRequest struct {
	context NitricContext
	payload []byte
}

// GetContext - return the context of a request, with metadata about that request.
func (n *NitricRequest) GetContext() NitricContext {
	return n.context
}

// GetPayload - return the []byte payload of the request.
func (n *NitricRequest) GetPayload() []byte {
	return n.payload
}

// GetStruct - Unmarshals the request body from JSON to the provided interface{}
func (n *NitricRequest) GetStruct(object interface{}) error {
	return json.Unmarshal(n.payload, object)
}

// contextFromHeaders - converts standard nitric HTTP headers into a context struct.
func contextFromHeaders(h fasthttp.RequestHeader) NitricContext {
	return NitricContext{
		requestID:   string(h.Peek("x-nitric-request-id")),
		sourceType:  sourceTypeFromString(string(h.Peek("x-nitric-source-type"))),
		source:      string(h.Peek("x-nitric-source")),
		payloadType: string(h.Peek("x-nitric-payload-type")),
	}
}

// fromHttpRequest - converts a standard nitric HTTP request into a NitricRequest to be passed to functions.
func fromRequestContext(ctx *fasthttp.RequestCtx) *NitricRequest {
	context := contextFromHeaders(ctx.Request.Header)

	return &NitricRequest{
		context: context,
		payload: ctx.Request.Body(),
	}
}
