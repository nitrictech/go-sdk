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

package events

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// Topic
type Topic interface {
	Name() string
	Publish(*Event) (*Event, error)
}

type topicImpl struct {
	name string
	ec   v1.EventServiceClient
}

func (s *topicImpl) Name() string {
	return s.name
}

func (s *topicImpl) Publish(evt *Event) (*Event, error) {
	// Convert payload to Protobuf Struct
	payloadStruct, err := structpb.NewStruct(evt.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %s", err)
	}

	r, err := s.ec.Publish(context.TODO(), &v1.EventPublishRequest{
		Topic: s.name,
		Event: &v1.NitricEvent{
			Id:          evt.ID,
			Payload:     payloadStruct,
			PayloadType: evt.PayloadType,
		},
	})

	if err != nil {
		return nil, err
	}

	// Return a reference to a new event with a populated ID
	return &Event{
		ID:          r.GetId(),
		Payload:     evt.Payload,
		PayloadType: evt.PayloadType,
	}, nil
}
