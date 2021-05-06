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

package eventclient

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type Event struct {
	Payload     map[string]interface{}
	PayloadType string
	ID          string
}

type EventClient interface {
	Publish(opts *PublishOptions) (*PublishResult, error)
}

type NitricEventClient struct {
	conn *grpc.ClientConn
	c    v1.EventClient
}

type PublishOptions struct {
	Topic string
	Event *Event
}

type PublishResult struct {
	RequestID string
}

// Publish - publishes the provided event data to the specified topic.
func (e NitricEventClient) Publish(opts *PublishOptions) (*PublishResult, error) {
	// Convert payload to Protobuf Struct
	payloadStruct, err := structpb.NewStruct(opts.Event.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %s", err)
	}

	// Publish the event
	resp, err := e.c.Publish(context.Background(), &v1.EventPublishRequest{
		Topic: opts.Topic,
		Event: &v1.NitricEvent{
			Id:          opts.Event.ID,
			PayloadType: opts.Event.PayloadType,
			Payload:     payloadStruct,
		},
	})

	if err != nil {
		return nil, err
	}

	return &PublishResult{
		RequestID: resp.Id,
	}, nil
}

func NewEventClient(conn *grpc.ClientConn) EventClient {
	return &NitricEventClient{
		conn: conn,
		c:    v1.NewEventClient(conn),
	}
}

func NewWithClient(eventClient v1.EventClient, topicClient v1.TopicClient) EventClient {
	return &NitricEventClient{
		c: eventClient,
	}
}
