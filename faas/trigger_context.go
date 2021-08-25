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
	pb "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
)

// NitricTriggerContext - Represents the contextual metadata for a Nitric function request.
type NitricTriggerContext struct {
	context interface{}
}

// IsHttp - Returns true if the trigger was via an HTTP Request and HTTP metadata is available
func (c *NitricTriggerContext) IsHttp() bool {
	_, ok := c.context.(*NitricHttpTriggerContext)

	return ok
}

// IsTopic - Returns true if the trigger was an Event via a Topic and Event/Topic metadata is available
func (c *NitricTriggerContext) IsTopic() bool {
	_, ok := c.context.(*NitricTopicTriggerContext)

	return ok
}

// AsHttp - Returns NitricTriggerContext as NitricHttpTriggerContext if IsHttp is true
// otherwise returns nil
func (c *NitricTriggerContext) AsHttp() *NitricHttpTriggerContext {
	if ctx, ok := c.context.(*NitricHttpTriggerContext); ok {
		return ctx
	}

	return nil
}

// AsTopic - Returns NitricTriggerContext as NitricTopicTriggerContext if IsTopic is true
// otherwise returns nil
func (c *NitricTriggerContext) AsTopic() *NitricTopicTriggerContext {
	if ctx, ok := c.context.(*NitricTopicTriggerContext); ok {
		return ctx
	}

	return nil
}

type NitricHttpTriggerContext struct {
	Method      string
	Headers     map[string][]string
	Path        string
	QueryParams map[string]string
}

type NitricTopicTriggerContext struct {
	Topic string
}

// ContextFromTriggerRequest - Creates TriggerContext from grpc TriggerRequest
func ContextFromTriggerRequest(grpcTrigger *pb.TriggerRequest) *NitricTriggerContext {
	triggerCtx := &NitricTriggerContext{}

	if grpcTrigger.GetHttp() != nil {
		http := grpcTrigger.GetHttp()

		headers := make(map[string][]string)
		if http.GetHeaders() != nil {
			for key, val := range http.GetHeaders() {
				headers[key] = val.Value
			}
		}

		triggerCtx.context = &NitricHttpTriggerContext{
			Method:      http.GetMethod(),
			Headers:     headers,
			Path:        http.GetPath(),
			QueryParams: http.GetQueryParams(),
		}

	} else if grpcTrigger.GetTopic() != nil {
		topic := grpcTrigger.GetTopic()
		triggerCtx.context = &NitricTopicTriggerContext{
			Topic: topic.Topic,
		}
	}

	// FIXME: Look at returning error over nil context
	return triggerCtx
}
