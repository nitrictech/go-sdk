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

package apis

import (
	"context"
	errorsstd "errors"

	grpcx "github.com/nitrictech/go-sdk/internal/grpc"
	"github.com/nitrictech/go-sdk/internal/handlers"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	"github.com/nitrictech/go-sdk/nitric/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
)

type apiWorker struct {
	client              v1.ApiClient
	Handler             handlers.Handler[Ctx]
	registrationRequest *v1.RegistrationRequest
}

type apiWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Handler             handlers.Handler[Ctx]
}

var _ workers.StreamWorker = (*apiWorker)(nil)

// Start runs the API worker, creating a stream to the Nitric server
func (a *apiWorker) Start(ctx context.Context) error {
	initReq := &v1.ClientMessage{
		Content: &v1.ClientMessage_RegistrationRequest{
			RegistrationRequest: a.registrationRequest,
		},
	}

	createStream := func(ctx context.Context) (workers.Stream[v1.ClientMessage, v1.RegistrationResponse, *v1.ServerMessage], error) {
		return a.client.Serve(ctx)
	}

	handlerSrvMsg := func(msg *v1.ServerMessage) (*v1.ClientMessage, error) {
		if msg.GetRegistrationResponse() != nil {
			// No need to respond to the registration response
			return nil, nil
		}

		if msg.GetHttpRequest() != nil {
			handlerCtx := NewCtx(msg)

			err := a.Handler(handlerCtx)
			if err != nil {
				handlerCtx.WithError(err)
			}

			return handlerCtx.ToClientMessage(), nil
		}

		return nil, errors.NewWithCause(
			codes.Internal,
			"ApiWorker: Unhandled server message",
			errorsstd.New("unhandled server message"),
		)
	}

	return workers.HandleStream(
		ctx,
		createStream,
		initReq,
		handlerSrvMsg,
	)
}

func newApiWorker(opts *apiWorkerOpts) *apiWorker {
	conn, err := grpcx.GetConnection()
	if err != nil {
		panic(errors.NewWithCause(
			codes.Unavailable,
			"NewApiWorker: Unable to reach ApiClient",
			err,
		))
	}

	client := v1.NewApiClient(conn)

	return &apiWorker{
		client:              client,
		registrationRequest: opts.RegistrationRequest,
		Handler:             opts.Handler,
	}
}
