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

	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

// Events
type Events interface {
	// Topic - Retrieve a Topic reference
	Topic(name string) Topic
	// Topics - Retrievs a list of available Topic references
	Topics() ([]Topic, error)
}

type eventsImpl struct {
	ec v1.EventClient
	tc v1.TopicClient
}

func (s *eventsImpl) Topic(name string) Topic {
	// Just return the straight topic reference
	// we can fail if the topic does not exist
	return &topicImpl{
		name: name,
		ec:   s.ec,
	}
}

func (s *eventsImpl) Topics() ([]Topic, error) {
	r, err := s.tc.List(context.TODO(), &v1.TopicListRequest{})

	if err != nil {
		return nil, err
	}

	ts := make([]Topic, 0)
	for _, topic := range r.GetTopics() {
		ts = append(ts, s.Topic(topic.GetName()))
	}

	return ts, nil
}

// New - Construct a new Eventing Client with default options
func New() (Events, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)

	if err != nil {
		return nil, err
	}

	ec := v1.NewEventClient(conn)
	tc := v1.NewTopicClient(conn)

	return &eventsImpl{
		ec: ec,
		tc: tc,
	}, nil
}
