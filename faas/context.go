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

	pb "github.com/nitrictech/apis/go/nitric/v1"
)

type TriggerContext interface {
	Http() *HttpContext
	Event() *EventContext
}

type triggerContextImpl struct {
	http  *HttpContext
	event *EventContext
}

func (t triggerContextImpl) Http() *HttpContext {
	return t.http
}

func (t triggerContextImpl) Event() *EventContext {
	return t.event
}

func triggerContextFromGrpcTriggerRequest(triggerReq *pb.TriggerRequest) (*triggerContextImpl, error) {
	trigCtx := &triggerContextImpl{}

	if triggerReq.GetHttp() != nil {
		httpTrig := triggerReq.GetHttp()

		headers := make(map[string][]string)
		if httpTrig.GetHeaders() != nil {
			for key, val := range httpTrig.GetHeaders() {
				headers[key] = val.Value
			}
		} else if httpTrig.GetHeadersOld() != nil {
			for key, val := range httpTrig.GetHeadersOld() {
				headers[key] = []string{val}
			}
		}

		query := make(map[string][]string)
		if httpTrig.GetQueryParams() != nil {
			for k, v := range httpTrig.GetQueryParams() {
				query[k] = v.Value
			}
		} else if httpTrig.GetQueryParamsOld() != nil {
			for k, v := range httpTrig.GetQueryParamsOld() {
				query[k] = []string{v}
			}
		}

		trigCtx.http = &HttpContext{
			Request: &httpRequestImpl{
				dataRequestImpl: dataRequestImpl{
					data:     triggerReq.GetData(),
					mimeType: triggerReq.GetMimeType(),
				},
				method:  httpTrig.GetMethod(),
				headers: headers,
				query:   query,
				path:    httpTrig.GetPath(),
			},
			Response: &HttpResponse{
				Status: 200,
				Headers: map[string][]string{
					"Content-Type": {"text/plain"},
				},
				Body: []byte("Success"),
			},
			extraContext: &extraContext{
				Extras: make(map[string]interface{}),
			},
		}
	} else if triggerReq.GetTopic() != nil {
		topic := triggerReq.GetTopic()

		trigCtx.event = &EventContext{
			Request: &eventRequestImpl{
				dataRequestImpl: dataRequestImpl{
					data:     triggerReq.GetData(),
					mimeType: triggerReq.GetMimeType(),
				},
				topic: topic.GetTopic(),
			},
			Response: &EventResponse{
				Success: true,
			},
			extraContext: &extraContext{
				Extras: make(map[string]interface{}),
			},
		}
	} else {
		return nil, fmt.Errorf("invalid trigger request")
	}

	return trigCtx, nil
}

func triggerContextToGrpcTriggerResponse(trig *triggerContextImpl) (*pb.TriggerResponse, error) {
	if trig.http != nil {

		headers := make(map[string]*pb.HeaderValue)
		headersOld := make(map[string]string)

		for k, v := range trig.http.Response.Headers {
			headersOld[k] = v[0]
			headers[k] = &pb.HeaderValue{
				Value: v,
			}
		}

		return &pb.TriggerResponse{
			Data: trig.http.Response.Body,
			Context: &pb.TriggerResponse_Http{
				Http: &pb.HttpResponseContext{
					Status:     int32(trig.http.Response.Status),
					Headers:    headers,
					HeadersOld: headersOld,
				},
			},
		}, nil
	} else if trig.event != nil {
		return &pb.TriggerResponse{
			// Don't actually need data available here
			Data: []byte(""),
			Context: &pb.TriggerResponse_Topic{
				Topic: &pb.TopicResponseContext{
					Success: trig.event.Response.Success,
				},
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported trigger context type")
}

type extraContext struct {
	Extras map[string]interface{}
}

type HttpContext struct {
	Request  HttpRequest
	Response *HttpResponse
	*extraContext
}

type EventContext struct {
	Request  EventRequest
	Response *EventResponse
	*extraContext
}
