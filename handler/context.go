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

package handler

import (
	http "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
)

type HttpContext struct {
	id       string
	Request  HttpRequest
	Response *HttpResponse
	Extras   map[string]interface{}
}

func (c *HttpContext) ToClientMessage() *http.ClientMessage {
	headers := make(map[string]*http.HeaderValue)
	for k, v := range c.Response.Headers {
		headers[k] = &http.HeaderValue{
			Value: v,
		}
	}

	return &http.ClientMessage{
		Id: c.id,
		Content: &http.ClientMessage_HttpResponse{
			HttpResponse: &http.HttpResponse{
				Status:  int32(c.Response.Status),
				Headers: headers,
				Body:    c.Response.Body,
			},
		},
	}
}

func NewHttpContext(msg *http.ServerMessage) *HttpContext {
	req := msg.GetHttpRequest()

	headers := make(map[string][]string)
	for k, v := range req.Headers {
		headers[k] = v.GetValue()
	}

	query := make(map[string][]string)
	for k, v := range req.QueryParams {
		query[k] = v.GetValue()
	}

	return &HttpContext{
		id: msg.Id,
		Request: &httpRequestImpl{
			method:     req.Method,
			path:       req.Path,
			pathParams: req.PathParams,
			query:      query,
			headers:    headers,
		},
		Response: &HttpResponse{
			Status:  200,
			Headers: map[string][]string{},
			Body:    nil,
		},
	}
}

func (c *HttpContext) WithError(err error) {
	c.Response = &HttpResponse{
		Status: 500,
		Headers: map[string][]string{
			"Content-Type": []string{"text/plain"},
		},
		Body: []byte(err.Error()),
	}
}

type MessageContext struct {
	Request  MessageRequest
	Response *MessageResponse
	Extras   map[string]interface{}
}

type IntervalContext struct {
	Request  IntervalRequest
	Response *IntervalResponse
	Extras   map[string]interface{}
}

type BlobEventContext struct {
	Request  BlobEventRequest
	Response *BlobEventResponse
	Extras   map[string]interface{}
}

type FileEventContext struct {
	Request  FileEventRequest
	Response *FileEventResponse
	Extras   map[string]interface{}
}

type WebsocketContext struct {
	Request  WebsocketRequest
	Response *WebsocketResponse
	Extras   map[string]interface{}
}