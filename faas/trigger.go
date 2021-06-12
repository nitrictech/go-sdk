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

	pb "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
)

// NitricRequest - represents a request to trigger a function, with payload and context required to execute that function.
type NitricTrigger struct {
	context *NitricTriggerContext
	data    []byte
}

// GetContext - return the context of a request, with metadata about that request.
func (n *NitricTrigger) GetContext() *NitricTriggerContext {
	return n.context
}

// GetPayload - return the []byte payload of the request.
func (n *NitricTrigger) GetData() []byte {
	return n.data
}

// GetStruct - Unmarshals the request body from JSON to the provided interface{}
func (n *NitricTrigger) GetDataAsStruct(object interface{}) error {
	return json.Unmarshal(n.data, object)
}

// Default Response - Returns a default response object dependent on the Trigger context
func (n *NitricTrigger) DefaultResponse() *NitricResponse {

	var context interface{} = nil

	if n.context.IsHttp() {
		context = &HttpResponseContext{
			Headers: make(map[string]string),
			Status:  200,
		}
	} else if n.context.IsHttp() {
		context = &TopicResponseContext{
			Success: true,
		}
	}

	return &NitricResponse{
		data: nil,
		context: &ResponseContext{
			context: context,
		},
	}
}

// fromHttpRequest - converts a standard nitric HTTP request into a NitricRequest to be passed to functions.
func FromGrpcTriggerRequest(triggerReq *pb.TriggerRequest) (*NitricTrigger, error) {
	context := ContextFromTriggerRequest(triggerReq)

	return &NitricTrigger{
		context: context,
		data:    triggerReq.GetData(),
	}, nil
}
