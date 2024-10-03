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
	"io"

	"google.golang.org/grpc"

	errorsstd "errors"

	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
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

// Start implements Worker.
func (s *jobWorker) Start(ctx context.Context) error {
	initReq := &v1.ClientMessage{
		Content: &v1.ClientMessage_RegistrationRequest{
			RegistrationRequest: s.registrationRequest,
		},
	}

	// Create the request stream and send the initial request
	stream, err := s.client.HandleJob(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(initReq)
	if err != nil {
		return err
	}
	for {
		var ctx *Ctx

		resp, err := stream.Recv()

		if errorsstd.Is(err, io.EOF) {
			err = stream.CloseSend()
			if err != nil {
				return err
			}

			return nil
		} else if err == nil && resp.GetRegistrationResponse() != nil {
			// Do nothing
		} else if err == nil && resp.GetJobRequest() != nil {
			ctx = NewCtx(resp)
			err = s.handler(ctx)
			if err != nil {
				ctx.WithError(err)
			}

			err = stream.Send(ctx.ToClientMessage())
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
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
