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

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

// Queue is a resource for async send/receive messaging.
type Queue interface {
	// Name - The name of the queue
	Name() string
	// Send - Push a number of tasks to a queue
	Send(context.Context, []*Task) ([]*FailedTask, error)
	// Receive - Retrieve tasks from a queue to a maximum of the given depth
	Receive(context.Context, int) ([]ReceivedTask, error)
}

type queueImpl struct {
	name string
	queueClient    v1.QueueServiceClient
}

func (q *queueImpl) Name() string {
	return q.name
}

func (q *queueImpl) Receive(ctx context.Context, depth int) ([]ReceivedTask, error) {
	if depth < 1 {
		return nil, errors.New(codes.InvalidArgument, "Queue.Receive: depth cannot be less than 1")
	}

	r, err := q.queueClient.Receive(ctx, &v1.QueueReceiveRequest{
		Queue: q.name,
		Depth: int32(depth),
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	rts := make([]ReceivedTask, len(r.GetTasks()))

	for i, task := range r.GetTasks() {
		rts[i] = &receivedTaskImpl{
			queue:   q.name,
			queueClient:      q.queueClient,
			leaseId: task.GetLeaseId(),
			task:    wireToTask(task),
		}
	}

	return rts, nil
}

func (q *queueImpl) Send(ctx context.Context, tasks []*Task) ([]*FailedTask, error) {
	// Convert SDK Task objects to gRPC Task objects
	wireTasks := make([]*v1.NitricTask, len(tasks))
	for i, task := range tasks {
		wireTask, err := taskToWire(task)
		if err != nil {
			return nil, errors.NewWithCause(
				codes.Internal,
				"Queue.Send: Unable to send tasks",
				err,
			)
		}
		wireTasks[i] = wireTask
	}

	// Push the tasks to the queue
	res, err := q.queueClient.SendBatch(ctx, &v1.QueueSendBatchRequest{
		Queue: q.name,
		Tasks: wireTasks,
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
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
