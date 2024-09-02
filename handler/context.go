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
	batch "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	schedules "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
	storage "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topics "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	websockets "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
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
			data:       req.Body,
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
			"Content-Type": {"text/plain"},
		},
		Body: []byte("Internal Server Error"),
	}
}

type MessageContext struct {
	id       string
	Request  MessageRequest
	Response *MessageResponse
	Extras   map[string]interface{}
}

func (c *MessageContext) ToClientMessage() *topics.ClientMessage {
	return &topics.ClientMessage{
		Id: c.id,
		Content: &topics.ClientMessage_MessageResponse{
			MessageResponse: &topics.MessageResponse{
				Success: true,
			},
		},
	}
}

func NewMessageContext(msg *topics.ServerMessage) *MessageContext {
	return &MessageContext{
		id: msg.Id,
		Request: &messageRequestImpl{
			topicName: msg.GetMessageRequest().TopicName,
			message:   msg.GetMessageRequest().Message.GetStructPayload().AsMap(),
		},
		Response: &MessageResponse{
			Success: true,
		},
	}
}

func (c *MessageContext) WithError(err error) {
	c.Response = &MessageResponse{
		Success: false,
	}
}

type JobContext struct {
	id       string
	Request  JobRequest
	Response *JobResponse
	Extras   map[string]interface{}
}

func (c *JobContext) ToClientMessage() *batch.ClientMessage {
	return &batch.ClientMessage{
		Id: c.id,
		Content: &batch.ClientMessage_JobResponse{
			JobResponse: &batch.JobResponse{
				Success: true,
			},
		},
	}
}

func NewJobContext(msg *batch.ServerMessage) *JobContext {
	return &JobContext{
		id: msg.Id,
		Request: &jobRequest{
			jobName: msg.GetJobRequest().JobName,
			data:    msg.GetJobRequest().Data.GetStruct().AsMap(),
		},
		Response: &JobResponse{
			Success: true,
		},
	}
}

func (c *JobContext) WithError(err error) {
	c.Response = &JobResponse{
		Success: false,
	}
}

type IntervalContext struct {
	id       string
	Request  IntervalRequest
	Response *IntervalResponse
	Extras   map[string]interface{}
}

func (c *IntervalContext) ToClientMessage() *schedules.ClientMessage {
	return &schedules.ClientMessage{
		Id: c.id,
		Content: &schedules.ClientMessage_IntervalResponse{
			IntervalResponse: &schedules.IntervalResponse{},
		},
	}
}

func NewIntervalContext(msg *schedules.ServerMessage) *IntervalContext {
	return &IntervalContext{
		id: msg.Id,
		Request: &intervalRequestImpl{
			scheduleName: msg.GetIntervalRequest().ScheduleName,
		},
		Response: &IntervalResponse{
			Success: true,
		},
	}
}

func (c *IntervalContext) WithError(err error) {
	c.Response = &IntervalResponse{
		Success: false,
	}
}

type BlobEventContext struct {
	id       string
	Request  BlobEventRequest
	Response *BlobEventResponse
	Extras   map[string]interface{}
}

func (c *BlobEventContext) ToClientMessage() *storage.ClientMessage {
	return &storage.ClientMessage{
		Id: c.id,
		Content: &storage.ClientMessage_BlobEventResponse{
			BlobEventResponse: &storage.BlobEventResponse{
				Success: c.Response.Success,
			},
		},
	}
}

func NewBlobEventContext(msg *storage.ServerMessage) *BlobEventContext {
	req := msg.GetBlobEventRequest()

	return &BlobEventContext{
		id: msg.Id,
		Request: &blobEventRequestImpl{
			key: req.GetBlobEvent().Key,
		},
		Response: &BlobEventResponse{
			Success: true,
		},
	}
}

func (c *BlobEventContext) WithError(err error) {
	c.Response = &BlobEventResponse{
		Success: false,
	}
}

type FileEventContext struct {
	Request  FileEventRequest
	Response *FileEventResponse
	Extras   map[string]interface{}
}

type WebsocketContext struct {
	id       string
	Request  WebsocketRequest
	Response *WebsocketResponse
	Extras   map[string]interface{}
}

func (c *WebsocketContext) ToClientMessage() *websockets.ClientMessage {
	return &websockets.ClientMessage{
		Id: c.id,
		Content: &websockets.ClientMessage_WebsocketEventResponse{
			WebsocketEventResponse: &websockets.WebsocketEventResponse{
				WebsocketResponse: &websockets.WebsocketEventResponse_ConnectionResponse{
					ConnectionResponse: &websockets.WebsocketConnectionResponse{
						Reject: c.Response.Reject,
					},
				},
			},
		},
	}
}

func NewWebsocketContext(msg *websockets.ServerMessage) *WebsocketContext {
	req := msg.GetWebsocketEventRequest()

	var eventType WebsocketEventType
	switch req.WebsocketEvent.(type) {
	case *websockets.WebsocketEventRequest_Disconnection:
		eventType = WebsocketDisconnect
	case *websockets.WebsocketEventRequest_Message:
		eventType = WebsocketMessage
	default:
		eventType = WebsocketConnect
	}

	queryParams := make(map[string]*http.QueryValue)
	if eventType == WebsocketConnect {
		for key, value := range req.GetConnection().GetQueryParams() {
			queryParams[key] = &http.QueryValue{
				Value: value.Value,
			}
		}
	}

	var _message string = ""
	if req.GetMessage() != nil {
		_message = string(req.GetMessage().Body)
	}

	return &WebsocketContext{
		id: msg.Id,
		Request: &websocketRequestImpl{
			socketName:   req.SocketName,
			eventType:    eventType,
			connectionId: req.ConnectionId,
			queryParams:  queryParams,
			message:      _message,
		},
		Response: &WebsocketResponse{
			Reject: false,
		},
	}
}

func (c *WebsocketContext) WithError(err error) {
	c.Response = &WebsocketResponse{
		Reject: true,
	}
}
