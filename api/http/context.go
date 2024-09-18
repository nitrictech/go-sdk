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

package http

import apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"

type Ctx struct {
	id       string
	Request  Request
	Response *Response
	Extras   map[string]interface{}
}

func (c *Ctx) ToClientMessage() *apispb.ClientMessage {
	headers := make(map[string]*apispb.HeaderValue)
	for k, v := range c.Response.Headers {
		headers[k] = &apispb.HeaderValue{
			Value: v,
		}
	}

	return &apispb.ClientMessage{
		Id: c.id,
		Content: &apispb.ClientMessage_HttpResponse{
			HttpResponse: &apispb.HttpResponse{
				Status:  int32(c.Response.Status),
				Headers: headers,
				Body:    c.Response.Body,
			},
		},
	}
}

func NewCtx(msg *apispb.ServerMessage) *Ctx {
	req := msg.GetHttpRequest()

	headers := make(map[string][]string)
	for k, v := range req.Headers {
		headers[k] = v.GetValue()
	}

	query := make(map[string][]string)
	for k, v := range req.QueryParams {
		query[k] = v.GetValue()
	}

	return &Ctx{
		id: msg.Id,
		Request: &RequestImpl{
			method:     req.Method,
			path:       req.Path,
			pathParams: req.PathParams,
			query:      query,
			headers:    headers,
			data:       req.Body,
		},
		Response: &Response{
			Status:  200,
			Headers: map[string][]string{},
			Body:    nil,
		},
	}
}

func (c *Ctx) WithError(err error) {
	c.Response = &Response{
		Status: 500,
		Headers: map[string][]string{
			"Content-Type": {"text/plain"},
		},
		Body: []byte("Internal Server Error"),
	}
}
