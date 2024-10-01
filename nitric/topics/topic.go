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

package topics

import (
	"fmt"

	"github.com/nitrictech/go-sdk/nitric/handlers"
	"github.com/nitrictech/go-sdk/nitric/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

// TopicPermission defines the available permissions on a topic
type TopicPermission string

const (
	// TopicPublish is required to call Publish on a topic.
	TopicPublish TopicPermission = "publish"
)

type SubscribableTopic interface {
	// Allow requests the given permissions to the topic.
	Allow(TopicPermission, ...TopicPermission) (*TopicClient, error)

	// Subscribe will register and start a subscription handler that will be called for all events from this topic.
	// Valid function signatures for handler are:
	//
	//	func()
	//	func() error
	//	func(*topics.Ctx)
	//	func(*topics.Ctx) error
	//	Handler[topics.Ctx]
	Subscribe(interface{})
}

type subscribableTopic struct {
	name         string
	manager      *workers.Manager
	registerChan <-chan workers.RegisterResult
}

// NewTopic creates a new Topic with the give name.
func NewTopic(name string) SubscribableTopic {
	topic := &subscribableTopic{
		name:    name,
		manager: workers.GetDefaultManager(),
	}

	topic.registerChan = topic.manager.RegisterResource(&v1.ResourceDeclareRequest{
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

func (t *subscribableTopic) Allow(permission TopicPermission, permissions ...TopicPermission) (*TopicClient, error) {
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

	err := t.manager.RegisterPolicy(registerResult.Identifier, actions...)
	if err != nil {
		return nil, err
	}

	return NewTopicClient(t.name)
}

func (t *subscribableTopic) Subscribe(handler interface{}) {
	registrationRequest := &topicspb.RegistrationRequest{
		TopicName: t.name,
	}

	typedHandler, err := handlers.HandlerFromInterface[Ctx](handler)
	if err != nil {
		panic(err)
	}

	opts := &subscriptionWorkerOpts{
		RegistrationRequest: registrationRequest,
		Handler:             typedHandler,
	}

	worker := newSubscriptionWorker(opts)
	t.manager.AddWorker("SubscriptionWorker:"+t.name, worker)
}
