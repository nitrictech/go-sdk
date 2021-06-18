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

// NitricResponse - represents the results of calling a function.
type NitricResponse struct {
	context *ResponseContext
	data    []byte
}

func (n *NitricResponse) SetData(data []byte) {
	n.data = data
}

func (n *NitricResponse) GetContext() *ResponseContext {
	return n.context
}

// ToTriggerResponse - Tranlates a Nitric Response for gRPC transport to the membrane
func (n *NitricResponse) ToTriggerResponse() *pb.TriggerResponse {

	triggerResponse := &pb.TriggerResponse{}

	triggerResponse.Data = n.data

	if n.context.IsHttp() {
		http := n.context.AsHttp()
		triggerResponse.Context = &pb.TriggerResponse_Http{
			Http: &pb.HttpResponseContext{
				Headers: http.Headers,
				Status:  int32(http.Status),
			},
		}
	} else if n.context.IsTopic() {
		topic := n.context.AsTopic()
		triggerResponse.Context = &pb.TriggerResponse_Topic{
			Topic: &pb.TopicResponseContext{
				Success: topic.Success,
			},
		}
	}

	return triggerResponse
}
