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

package queues

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
)

// Queue - Interface for a Queue reference
type Queue interface {
	// Name - The name of the queue
	Name() string
	// Send - Push a number of tasks to a queue
	Send([]*Task) ([]*FailedTask, error)
	// Receive - Retrieve tasks from a queue to a maximum of the given depth
	Receive(int) ([]ReceivedTask, error)
}

type queueImpl struct {
	name string
	c    v1.QueueServiceClient
}

func (q *queueImpl) Name() string {
	return q.name
}

func (q *queueImpl) Receive(depth int) ([]ReceivedTask, error) {
	if depth < 1 {
		return nil, fmt.Errorf("Depth cannot be less than 1")
	}

	r, err := q.c.Receive(context.TODO(), &v1.QueueReceiveRequest{
		Queue: q.name,
		Depth: int32(depth),
	})

	if err != nil {
		return nil, err
	}

	rts := make([]ReceivedTask, len(r.GetTasks()))

	for i, task := range r.GetTasks() {
		rts[i] = &receivedTaskImpl{
			queue:   q.name,
			qc:      q.c,
			leaseId: task.GetLeaseId(),
			task:    wireToTask(task),
		}
	}

	return rts, nil
}

func (q *queueImpl) Send(tasks []*Task) ([]*FailedTask, error) {
	// Convert SDK Task objects to gRPC Task objects
	wireTasks := make([]*v1.NitricTask, len(tasks))
	for i, task := range tasks {
		wireTask, err := taskToWire(task)
		if err != nil {
			return nil, err
		}
		wireTasks[i] = wireTask
	}

	// Push the tasks to the queue
	res, err := q.c.SendBatch(context.Background(), &v1.QueueSendBatchRequest{
		Queue: q.name,
		Tasks: wireTasks,
	})
	if err != nil {
		return nil, err
	}

	// Convert the gRPC Failed Tasks to SDK Failed Task objects
	failedTasks := make([]*FailedTask, len(res.GetFailedTasks()))
	for i, failedTask := range res.GetFailedTasks() {
		failedTasks[i] = &FailedTask{
			Task:   wireToTask(failedTask.GetTask()),
			Reason: failedTask.GetMessage(),
		}
	}

	return failedTasks, nil
}
