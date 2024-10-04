// Copyright 2023 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package batch

import (
	"context"

	"google.golang.org/grpc"

	errorsstd "errors"

	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	"github.com/nitrictech/go-sdk/nitric/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
)

type jobWorker struct {
	client              v1.JobClient
	registrationRequest *v1.RegistrationRequest
	handler             Handler
}
type jobWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Handler             Handler
}

// Start runs the Job worker, creating a stream to the Nitric server
func (s *jobWorker) Start(ctx context.Context) error {
	initReq := &v1.ClientMessage{
		Content: &v1.ClientMessage_RegistrationRequest{
			RegistrationRequest: s.registrationRequest,
		},
	}

	createStream := func(ctx context.Context) (workers.Stream[v1.ClientMessage, v1.RegistrationResponse, *v1.ServerMessage], error) {
		return s.client.HandleJob(ctx)
	}

	handleSrvMsg := func(msg *v1.ServerMessage) (*v1.ClientMessage, error) {
		if msg.GetJobRequest() != nil {
			handlerCtx := NewCtx(msg)

			err := s.handler(handlerCtx)
			if err != nil {
				handlerCtx.WithError(err)
			}

			return handlerCtx.ToClientMessage(), nil
		}

		return nil, errors.NewWithCause(
			codes.Internal,
			"JobWorker: Unhandled server message",
			errorsstd.New("unhandled server message"),
		)
	}

	return workers.HandleStream(ctx, createStream, initReq, handleSrvMsg)
}

func newJobWorker(opts *jobWorkerOpts) *jobWorker {
	conn, err := grpc.NewClient(constants.NitricAddress(), constants.DefaultOptions()...)
	if err != nil {
		panic(errors.NewWithCause(
			codes.Unavailable,
			"NewJobWorker: Unable to reach JobClient",
			err,
		))
	}

	client := v1.NewJobClient(conn)

	return &jobWorker{
		client:              client,
		registrationRequest: opts.RegistrationRequest,
		handler:             opts.Handler,
	}
}
