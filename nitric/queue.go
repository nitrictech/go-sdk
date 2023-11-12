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
	nitricv1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

type QueuePermission string

const (
	QueueSending   QueuePermission = "sending"
	QueueReceiving QueuePermission = "receiving"
)

var QueueEverything []QueuePermission = []QueuePermission{QueueSending, QueueReceiving}

type Queue interface {
	With(QueuePermission, ...QueuePermission) (queues.Queue, error)
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
func (q *queue) With(permission QueuePermission, permissions ...QueuePermission) (queues.Queue, error) {
	allPerms := append([]QueuePermission{permission}, permissions...)

	return defaultManager.newQueue(q.name, allPerms...)
}

func (m *manager) newQueue(name string, permissions ...QueuePermission) (queues.Queue, error) {
	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	colRes := &nitricv1.Resource{
		Type: nitricv1.ResourceType_Queue,
		Name: name,
	}

	dr := &nitricv1.ResourceDeclareRequest{
		Resource: colRes,
		Config: &nitricv1.ResourceDeclareRequest_Queue{
			Queue: &nitricv1.QueueResource{},
		},
	}
	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []nitricv1.Action{}
	for _, perm := range permissions {
		switch perm {
		case QueueReceiving:
			actions = append(actions, nitricv1.Action_QueueReceive)
		case QueueSending:
			actions = append(actions, nitricv1.Action_QueueSend)
		default:
			return nil, fmt.Errorf("QueuePermission %s unknown", perm)
		}
	}
	if len(actions) > 0 {
		actions = append(actions, nitricv1.Action_QueueDetail, nitricv1.Action_QueueList)
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
