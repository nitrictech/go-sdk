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

	v1 "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

type Task struct {
	// ID - Unique ID for the task
	ID string
	// PayloadType - Deserialization hint for interprocess communication
	PayloadType string
	// Payload - The payload to include in this task
	Payload map[string]interface{}
}

type ReceivedTask interface {
	// Queue - Returns the name of the queue this task was retrieved from
	Queue() string
	// Task - Returns the Task data contained in this Recieved Task instance
	Task() *Task
	// Complete - Completes the task removing it from the queue
	Complete() error
}

type receivedTaskImpl struct {
	queue   string
	qc      v1.QueueServiceClient
	leaseId string
	task    *Task
}

func (r *receivedTaskImpl) Task() *Task {
	return r.task
}

func (r *receivedTaskImpl) Queue() string {
	return r.queue
}

func (r *receivedTaskImpl) Complete() error {
	_, err := r.qc.Complete(context.TODO(), &v1.QueueCompleteRequest{
		Queue:   r.queue,
		LeaseId: r.leaseId,
	})

	return err
}

type FailedTask struct {
	// Task - The task that failed to queue
	Task *Task
	// Reason - Reason for the failure
	Reason string
}

func taskToWire(task *Task) (*v1.NitricTask, error) {
	// Convert payload to Protobuf Struct
	payloadStruct, err := structpb.NewStruct(task.Payload)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Internal,
			"taskToWire: failed to serialize payload: %s",
			err,
		)
	}

	return &v1.NitricTask{
		Id:          task.ID,
		PayloadType: task.PayloadType,
		Payload:     payloadStruct,
	}, nil
}

func wireToTask(task *v1.NitricTask) *Task {
	return &Task{
		ID:          task.GetId(),
		PayloadType: task.GetPayloadType(),
		Payload:     task.GetPayload().AsMap(),
	}
}
