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

package resources

import (
	"context"
	"fmt"

	nitricv1 "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/events"
	"github.com/nitrictech/go-sdk/faas"
)

type TopicPermission string

const (
	TopicPublishing TopicPermission = "publishing"
)

type topic struct {
	events.Topic
	mgr *manager
}

type Topic interface {
	events.Topic
	Subscribe(...faas.EventMiddleware)
}

func NewTopic(name string, permissions ...TopicPermission) (Topic, error) {
	return run.NewTopic(name, permissions...)
}

func (m *manager) NewTopic(name string, permissions ...TopicPermission) (Topic, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	res := &nitricv1.Resource{
		Type: nitricv1.ResourceType_Topic,
		Name: name,
	}

	dr := &nitricv1.ResourceDeclareRequest{
		Resource: res,
		Config: &nitricv1.ResourceDeclareRequest_Topic{
			Topic: &nitricv1.TopicResource{},
		}}
	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []nitricv1.Action{}
	for _, perm := range permissions {
		switch perm {
		case TopicPublishing:
			actions = append(actions, nitricv1.Action_TopicDetail, nitricv1.Action_TopicEventPublish)
		default:
			return nil, fmt.Errorf("TopicPermission %s unknown", perm)
		}
	}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(res, actions))
	if err != nil {
		return nil, err
	}

	if m.evts == nil {
		evts, err := events.New()
		if err != nil {
			return nil, err
		}
		m.evts = evts
	}

	return &topic{
		Topic: m.evts.Topic(name),
		mgr:   m,
	}, nil
}

func (t *topic) Subscribe(middleware ...faas.EventMiddleware) {
	f := faas.New()
	f.Event(middleware...)
	f.WithSubscriptionWorkerOpts(faas.SubscriptionWorkerOptions{Topic: t.Name()})

	t.mgr.addStarter(fmt.Sprintf("topic:subscribe %s", t.Name()), f)
}
