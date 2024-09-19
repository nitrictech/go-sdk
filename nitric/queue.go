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

	"github.com/nitrictech/go-sdk/nitric/queues"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type QueuePermission string

const (
	QueueEnqueue QueuePermission = "enqueue"
	QueueDequeue QueuePermission = "dequeue"
)

var QueueEverything []QueuePermission = []QueuePermission{QueueEnqueue, QueueDequeue}

type Queue interface {
	// Allow requests the given permissions to the queue.
	Allow(QueuePermission, ...QueuePermission) (*queues.QueueClient, error)
}

type queue struct {
	name         string
	manager      *manager
	registerChan <-chan RegisterResult
}

// NewQueue - Create a new Queue resource
func NewQueue(name string) *queue {
	queue := &queue{
		name:         name,
		manager:      defaultManager,
		registerChan: make(chan RegisterResult),
	}

	queue.registerChan = defaultManager.registerResource(&v1.ResourceDeclareRequest{
		Id: &v1.ResourceIdentifier{
			Type: v1.ResourceType_Queue,
			Name: name,
		},
		Config: &v1.ResourceDeclareRequest_Queue{
			Queue: &v1.QueueResource{},
		},
	})

	return queue
}

func (q *queue) Allow(permission QueuePermission, permissions ...QueuePermission) (*queues.QueueClient, error) {
	allPerms := append([]QueuePermission{permission}, permissions...)

	actions := []v1.Action{}
	for _, perm := range allPerms {
		switch perm {
		case QueueDequeue:
			actions = append(actions, v1.Action_QueueDequeue)
		case QueueEnqueue:
			actions = append(actions, v1.Action_QueueEnqueue)
		default:
			return nil, fmt.Errorf("QueuePermission %s unknown", perm)
		}
	}

	registerResult := <-q.registerChan
	if registerResult.Err != nil {
		return nil, registerResult.Err
	}

	err := q.manager.registerPolicy(registerResult.Identifier, actions...)
	if err != nil {
		return nil, err
	}

	return queues.NewQueueClient(q.name)
}
