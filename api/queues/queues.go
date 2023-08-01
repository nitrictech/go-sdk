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
	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

// Queues - Idiomatic interface for the nitric queue service
type Queues interface {
	Queue(string) Queue
}

type queuesImpl struct {
	queueClient v1.QueueServiceClient
}

func (q *queuesImpl) Queue(name string) Queue {
	return &queueImpl{
		name: name,
		queueClient:    q.queueClient,
	}
}

// New - Construct a new Queueing Client with default options
func New() (Queues, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"Queues.New: Unable to reach QueueServiceServer",
			err,
		)
	}

	qClient := v1.NewQueueServiceClient(conn)

	return &queuesImpl{
		queueClient: qClient,
	}, nil
}
