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

package nitric

import (
	"context"
	"fmt"

	"github.com/nitrictech/go-sdk/api/topics"
	"github.com/nitrictech/go-sdk/handler"

	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

// TopicPermission defines the available permissions on a topic
type TopicPermission string

const (
	// TopicPublishing is required to call Publish on a topic.
	TopicPublish TopicPermission = "publish"
)

type Topic interface {
	topics.Topic
}

type SubscribableTopic interface {
	Allow(TopicPermission, ...TopicPermission) (Topic, error)

	// Subscribe will register and start a subscription handler that will be called for all events from this topic.
	Subscribe(...handler.MessageMiddleware)
}

type topic struct {
	topics.Topic

	manager Manager
}

type subscribableTopic struct {
	name    string
	manager Manager
}

// NewTopic creates a new Topic with the give permissions.
func NewTopic(name string) SubscribableTopic {
	return &subscribableTopic{
		name:    name,
		manager: defaultManager,
	}
}

func (t *subscribableTopic) Allow(permission TopicPermission, permissions ...TopicPermission) (Topic, error) {
	allPerms := append([]TopicPermission{permission}, permissions...)

	return defaultManager.newTopic(t.name, allPerms...)
}

func (m *manager) newTopic(name string, permissions ...TopicPermission) (Topic, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	res := &v1.ResourceIdentifier{
		Type: v1.ResourceType_Topic,
		Name: name,
	}

	dr := &v1.ResourceDeclareRequest{
		Id: res,
		Config: &v1.ResourceDeclareRequest_Topic{
			Topic: &v1.TopicResource{},
		},
	}
	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []v1.Action{}
	for _, perm := range permissions {
		switch perm {
		case TopicPublish:
			actions = append(actions, v1.Action_TopicPublish)
		default:
			return nil, fmt.Errorf("TopicPermission %s unknown", perm)
		}
	}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(res, actions))
	if err != nil {
		return nil, err
	}

	if m.topics == nil {
		evts, err := topics.New()
		if err != nil {
			return nil, err
		}
		m.topics = evts
	}

	return &topic{
		Topic:   m.topics.Topic(name),
		manager: m,
	}, nil
}

func (t *subscribableTopic) Subscribe(middleware ...handler.MessageMiddleware) {
	// TODO: create subscription worker
}
