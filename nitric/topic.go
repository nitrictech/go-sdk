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
	"fmt"

	"github.com/nitrictech/go-sdk/api/topics"

	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
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
	// Valid function signatures for middleware are:
	//
	//	func()
	//	func() error
	//	func(*topics.Ctx)
	//	func(*topics.Ctx) error
	//	func(*topics.Ctx) *topics.Ctx
	//	func(*topics.Ctx) (*topics.Ctx, error)
	//	func(*topics.Ctx, Handler[topics.Ctx]) *topics.Ctx
	//	func(*topics.Ctx, Handler[topics.Ctx]) error
	//	func(*topics.Ctx, Handler[topics.Ctx]) (*topics.Ctx, error)
	//	Middleware[topics.Ctx]
	//	Handler[topics.Ctx]
	Subscribe(...interface{})
}

type topic struct {
	topics.Topic

	manager Manager
}

type subscribableTopic struct {
	name         string
	manager      Manager
	registerChan <-chan RegisterResult
}

// NewTopic creates a new Topic with the give permissions.
func NewTopic(name string) SubscribableTopic {
	topic := &subscribableTopic{
		name:    name,
		manager: defaultManager,
	}

	topic.registerChan = defaultManager.registerResource(&v1.ResourceDeclareRequest{
		Id: &v1.ResourceIdentifier{
			Type: v1.ResourceType_Topic,
			Name: name,
		},
		Config: &v1.ResourceDeclareRequest_Topic{
			Topic: &v1.TopicResource{},
		},
	})

	return topic
}

func (t *subscribableTopic) Allow(permission TopicPermission, permissions ...TopicPermission) (Topic, error) {
	allPerms := append([]TopicPermission{permission}, permissions...)

	actions := []v1.Action{}
	for _, perm := range allPerms {
		switch perm {
		case TopicPublish:
			actions = append(actions, v1.Action_TopicPublish)
		default:
			return nil, fmt.Errorf("TopicPermission %s unknown", perm)
		}
	}

	registerResult := <-t.registerChan
	if registerResult.Err != nil {
		return nil, registerResult.Err
	}

	m, err := t.manager.registerPolicy(registerResult.Identifier, actions...)
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
		Topic:   m.topics.Topic(t.name),
		manager: m,
	}, nil
}

func (t *subscribableTopic) Subscribe(middleware ...interface{}) {
	registrationRequest := &topicspb.RegistrationRequest{
		TopicName: t.name,
	}

	middlewares, err := interfacesToMiddleware[topics.Ctx](middleware)
	if err != nil {
		panic(err)
	}

	composeHandler := Compose(middlewares...)

	opts := &subscriptionWorkerOpts{
		RegistrationRequest: registrationRequest,
		Middleware:          composeHandler,
	}

	worker := newSubscriptionWorker(opts)
	t.manager.addWorker("SubscriptionWorker:"+t.name, worker)
}
