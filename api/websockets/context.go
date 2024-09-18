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

package websockets

import websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"

type Ctx struct {
	id       string
	Request  Request
	Response *Response
	Extras   map[string]interface{}
}

func (c *Ctx) ToClientMessage() *websocketspb.ClientMessage {
	return &websocketspb.ClientMessage{
		Id: c.id,
		Content: &websocketspb.ClientMessage_WebsocketEventResponse{
			WebsocketEventResponse: &websocketspb.WebsocketEventResponse{
				WebsocketResponse: &websocketspb.WebsocketEventResponse_ConnectionResponse{
					ConnectionResponse: &websocketspb.WebsocketConnectionResponse{
						Reject: c.Response.Reject,
					},
				},
			},
		},
	}
}

func NewCtx(msg *websocketspb.ServerMessage) *Ctx {
	req := msg.GetWebsocketEventRequest()

	var eventType EventType
	switch req.WebsocketEvent.(type) {
	case *websocketspb.WebsocketEventRequest_Disconnection:
		eventType = EventType_Disconnect
	case *websocketspb.WebsocketEventRequest_Message:
		eventType = EventType_Message
	default:
		eventType = EventType_Connect
	}

	queryParams := make(map[string]*websocketspb.QueryValue)
	if eventType == EventType_Connect {
		for key, value := range req.GetConnection().GetQueryParams() {
			queryParams[key] = &websocketspb.QueryValue{
				Value: value.Value,
			}
		}
	}

	var _message string = ""
	if req.GetMessage() != nil {
		_message = string(req.GetMessage().Body)
	}

	return &Ctx{
		id: msg.Id,
		Request: &requestImpl{
			socketName:   req.SocketName,
			eventType:    eventType,
			connectionId: req.ConnectionId,
			queryParams:  queryParams,
			message:      _message,
		},
		Response: &Response{
			Reject: false,
		},
	}
}

func (c *Ctx) WithError(err error) {
	c.Response = &Response{
		Reject: true,
	}
}
