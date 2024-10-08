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

	grpcx "github.com/nitrictech/go-sdk/internal/grpc"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
)

// QueueClientIface is a resource for async enqueueing/dequeueing of messages.
type QueueClientIface interface {
	// Name - The name of the queue
	Name() string
	// Enqueue - Push a number of messages to a queue
	Enqueue(ctx context.Context, messages []map[string]interface{}) ([]*FailedMessage, error)
	// Dequeue - Retrieve messages from a queue to a maximum of the given depth
	Dequeue(ctx context.Context, depth int) ([]ReceivedMessage, error)
}

type QueueClient struct {
	name        string
	queueClient v1.QueuesClient
}

func (q *QueueClient) Name() string {
	return q.name
}

func (q *QueueClient) Dequeue(ctx context.Context, depth int) ([]ReceivedMessage, error) {
	if depth < 1 {
		return nil, errors.New(codes.InvalidArgument, "Queue.Dequeue: depth cannot be less than 1")
	}

	r, err := q.queueClient.Dequeue(ctx, &v1.QueueDequeueRequest{
		QueueName: q.name,
		Depth:     int32(depth),
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	rts := make([]ReceivedMessage, len(r.GetMessages()))

	for i, message := range r.GetMessages() {
		rts[i] = &leasedMessage{
			queueName:   q.name,
			queueClient: q.queueClient,
			leaseId:     message.GetLeaseId(),
			message:     wireToMessage(message.GetMessage()),
		}
	}

	return rts, nil
}

func (q *QueueClient) Enqueue(ctx context.Context, messages []map[string]interface{}) ([]*FailedMessage, error) {
	// Convert SDK Message objects to gRPC Message objects
	wireMessages := make([]*v1.QueueMessage, len(messages))
	for i, message := range messages {
		wireMessage, err := messageToWire(message)
		if err != nil {
			return nil, errors.NewWithCause(
				codes.Internal,
				"Queue.Enqueue: Unable to enqueue messages",
				err,
			)
		}
		wireMessages[i] = wireMessage
	}

	// Push the messages to the queue
	res, err := q.queueClient.Enqueue(ctx, &v1.QueueEnqueueRequest{
		QueueName: q.name,
		Messages:  wireMessages,
	})
	if err != nil {
		return nil, errors.FromGrpcError(err)
	}

	// Convert the gRPC Failed Messages to SDK Failed Message objects
	failedMessages := make([]*FailedMessage, len(res.GetFailedMessages()))
	for i, failedMessage := range res.GetFailedMessages() {
		failedMessages[i] = &FailedMessage{
			Message: wireToMessage(failedMessage.GetMessage()),
			Reason:  failedMessage.GetDetails(),
		}
	}

	return failedMessages, nil
}

func NewQueueClient(name string) (*QueueClient, error) {
	conn, err := grpcx.GetConnection()
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"NewQueueClient: unable to reach nitric server",
			err,
		)
	}

	queueClient := v1.NewQueuesClient(conn)

	return &QueueClient{
		name:        name,
		queueClient: queueClient,
	}, nil
}
