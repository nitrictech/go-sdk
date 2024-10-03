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

	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
	"github.com/nitrictech/protoutils"
)

type ReceivedMessage interface {
	// Queue - Returns the name of the queue this message was retrieved from
	Queue() string
	// Message - Returns the Message data contained in this Received Message instance
	Message() map[string]interface{}
	// Complete - Completes the message removing it from the queue
	Complete(context.Context) error
}

type leasedMessage struct {
	queueName   string
	queueClient v1.QueuesClient
	leaseId     string
	message     map[string]interface{}
}

func (r *leasedMessage) Message() map[string]interface{} {
	return r.message
}

func (r *leasedMessage) Queue() string {
	return r.queueName
}

func (r *leasedMessage) Complete(ctx context.Context) error {
	_, err := r.queueClient.Complete(ctx, &v1.QueueCompleteRequest{
		QueueName: r.queueName,
		LeaseId:   r.leaseId,
	})

	return err
}

type FailedMessage struct {
	// Message - The message that failed to queue
	Message map[string]interface{}
	// Reason - Reason for the failure
	Reason string
}

func messageToWire(message map[string]interface{}) (*v1.QueueMessage, error) {
	// Convert payload to Protobuf Struct
	payloadStruct, err := protoutils.NewStruct(message)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Internal,
			"messageToWire: failed to serialize message: %s",
			err,
		)
	}

	return &v1.QueueMessage{
		Content: &v1.QueueMessage_StructPayload{
			StructPayload: payloadStruct,
		},
	}, nil
}

func wireToMessage(message *v1.QueueMessage) map[string]interface{} {
	// TODO: verify that AsMap() ignores the proto field values
	return message.GetStructPayload().AsMap()
}
