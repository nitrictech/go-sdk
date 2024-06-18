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

	"github.com/nitrictech/go-sdk/api/queues"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type QueuePermission string

const (
	QueueEnqueue QueuePermission = "enqueue"
	QueueDequeue QueuePermission = "dequeue"
)

var QueueEverything []QueuePermission = []QueuePermission{QueueEnqueue, QueueDequeue}

type Queue interface {
	Allow(QueuePermission, ...QueuePermission) (queues.Queue, error)
}

type queue struct {
	name    string
	manager Manager
}

func NewQueue(name string) *queue {
	return &queue{
		name:    name,
		manager: defaultManager,
	}
}

// NewQueue registers this queue as a required resource for the calling function/container.
func (q *queue) Allow(permission QueuePermission, permissions ...QueuePermission) (queues.Queue, error) {
	allPerms := append([]QueuePermission{permission}, permissions...)

	return defaultManager.newQueue(q.name, allPerms...)
}

func (m *manager) newQueue(name string, permissions ...QueuePermission) (queues.Queue, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	colRes := &v1.ResourceIdentifier{
		Type: v1.ResourceType_Queue,
		Name: name,
	}

	dr := &v1.ResourceDeclareRequest{
		Id: colRes,
		Config: &v1.ResourceDeclareRequest_Queue{
			Queue: &v1.QueueResource{},
		},
	}
	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []v1.Action{}
	for _, perm := range permissions {
		switch perm {
		case QueueDequeue:
			actions = append(actions, v1.Action_QueueDequeue)
		case QueueEnqueue:
			actions = append(actions, v1.Action_QueueEnqueue)
		default:
			return nil, fmt.Errorf("QueuePermission %s unknown", perm)
		}
	}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(colRes, actions))
	if err != nil {
		return nil, err
	}

	if m.queues == nil {
		m.queues, err = queues.New()
		if err != nil {
			return nil, err
		}
	}

	return m.queues.Queue(name), nil
}
