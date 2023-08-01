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
	"fmt"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

type TriggerContext interface {
	Http() *HttpContext
	Event() *EventContext
	BucketNotification() *BucketNotificationContext
	Websocket() *WebsocketContext
}

type triggerContextImpl struct {
	http               *HttpContext
	event              *EventContext
	bucketNotification *BucketNotificationContext
	websocket *WebsocketContext
}

func (t triggerContextImpl) Http() *HttpContext {
	return t.http
}

func (t triggerContextImpl) Event() *EventContext {
	return t.event
}

func (t triggerContextImpl) BucketNotification() *BucketNotificationContext {
	return t.bucketNotification
}

func (t triggerContextImpl) Websocket() *WebsocketContext {
	return t.websocket
}

func triggerContextFromGrpcTriggerRequest(triggerReq *v1.TriggerRequest) (*triggerContextImpl, error) {
	trigCtx := &triggerContextImpl{}

	tc := map[string]string{}
	if triggerReq.TraceContext != nil {
		tc = triggerReq.TraceContext.GetValues()
	}

	if triggerReq.GetHttp() != nil {
		httpTrig := triggerReq.GetHttp()

		headers := make(map[string][]string)
		if httpTrig.GetHeaders() != nil {
			for key, val := range httpTrig.GetHeaders() {
				headers[key] = val.Value
			}
		}

		query := make(map[string][]string)
		if httpTrig.GetQueryParams() != nil {
			for k, v := range httpTrig.GetQueryParams() {
				query[k] = v.Value
			}
		}

		trigCtx.http = &HttpContext{
			Request: &httpRequestImpl{
				dataRequestImpl: dataRequestImpl{
					data:         triggerReq.GetData(),
					mimeType:     triggerReq.GetMimeType(),
					traceContext: tc,
				},
				method:     httpTrig.GetMethod(),
				headers:    headers,
				query:      query,
				pathParams: httpTrig.PathParams,
				path:       httpTrig.GetPath(),
			},
			Response: &HttpResponse{
				Status: 200,
				Headers: map[string][]string{
					"Content-Type": {"text/plain"},
				},
				Body: []byte("Success"),
			},
			Extras: make(map[string]interface{}),
		}
	} else if triggerReq.GetTopic() != nil {
		topic := triggerReq.GetTopic()

		trigCtx.event = &EventContext{
			Request: &eventRequestImpl{
				dataRequestImpl: dataRequestImpl{
					data:         triggerReq.GetData(),
					mimeType:     triggerReq.GetMimeType(),
					traceContext: tc,
				},
				topic: topic.GetTopic(),
			},
			Response: &EventResponse{
				Success: true,
			},
			Extras: make(map[string]interface{}),
		}
	} else if triggerReq.GetWebsocket() != nil {
		websocket := triggerReq.GetWebsocket()

		var evtType WebsocketEventType

		switch websocket.Event {
		case v1.WebsocketEvent_Connect:
			evtType = WebsocketConnect
		case v1.WebsocketEvent_Disconnect:
			evtType = WebsocketDisconnect
		case v1.WebsocketEvent_Message:
			evtType = WebsocketMessage
		}

		trigCtx.websocket = &WebsocketContext{
			Request: &websocketRequestImpl{
				dataRequestImpl: dataRequestImpl{
					data:         triggerReq.GetData(),
					mimeType:     triggerReq.GetMimeType(),
					traceContext: tc,
				},
				socket: websocket.Socket,
				connectionId: websocket.ConnectionId,
				eventType: evtType,
				queryParams: websocket.QueryParams,				
			},
			Response: &WebsocketResponse{
				Success: true,
			},
			Extras: make(map[string]interface{}),
		}
	} else if triggerReq.GetNotification() != nil {
		notification := triggerReq.GetNotification()

		var notificationType NotificationType

		switch notification.GetBucket().Type {
		case v1.BucketNotificationType_Created:
			notificationType = WriteNotification
		case v1.BucketNotificationType_Deleted:
			notificationType = DeleteNotification
		default:
			return nil, fmt.Errorf("notification type %s is not supported", notification.GetBucket().Type)
		}

		trigCtx.bucketNotification = &BucketNotificationContext{
			Request: &bucketNotificationRequestImpl{
				key: notification.GetBucket().Key,
				notificationType: notificationType,
			},
			Response: &BucketNotificationResponse{
				Success: true,
			},
			Extras: make(map[string]interface{}),
		}
	} else {
		return nil, fmt.Errorf("invalid trigger request")
	}

	return trigCtx, nil
}

func triggerContextToGrpcTriggerResponse(trig *triggerContextImpl) (*v1.TriggerResponse, error) {
	if trig.http != nil {
		headers := make(map[string]*v1.HeaderValue)
		headersOld := make(map[string]string)

		for k, v := range trig.http.Response.Headers {
			headersOld[k] = v[0]
			headers[k] = &v1.HeaderValue{
				Value: v,
			}
		}

		return &v1.TriggerResponse{
			Data: trig.http.Response.Body,
			Context: &v1.TriggerResponse_Http{
				Http: &v1.HttpResponseContext{
					Status:     int32(trig.http.Response.Status),
					Headers:    headers,
					HeadersOld: headersOld,
				},
			},
		}, nil
	} else if trig.event != nil {
		return &v1.TriggerResponse{
			// Don't actually need data available here
			Data: []byte(""),
			Context: &v1.TriggerResponse_Topic{
				Topic: &v1.TopicResponseContext{
					Success: trig.event.Response.Success,
				},
			},
		}, nil
	} else if trig.websocket != nil {
		return &v1.TriggerResponse{
			Data: []byte(""),
			Context: &v1.TriggerResponse_Websocket{
				Websocket: &v1.WebsocketResponseContext{
					Success: trig.websocket.Response.Success,
				},
			},
		}, nil
	} else if trig.bucketNotification != nil {
		return &v1.TriggerResponse{
			Data: []byte(""),
			Context: &v1.TriggerResponse_Notification{
				Notification: &v1.NotificationResponseContext{
					Success: trig.bucketNotification.Response.Success,
				},
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported trigger context type")
}

type HttpContext struct {
	Request  HttpRequest
	Response *HttpResponse
	Extras   map[string]interface{}
}

type EventContext struct {
	Request  EventRequest
	Response *EventResponse
	Extras   map[string]interface{}
}

type BucketNotificationContext struct {
	Request  BucketNotificationRequest
	Response *BucketNotificationResponse
	Extras   map[string]interface{}
}

type WebsocketContext struct {
	Request WebsocketRequest
	Response *WebsocketResponse
	Extras map[string]interface{}
}