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

package websockets

import (
	"context"
	errorsstd "errors"

	grpcx "github.com/nitrictech/go-sdk/internal/grpc"
	"github.com/nitrictech/go-sdk/internal/handlers"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	"github.com/nitrictech/go-sdk/nitric/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
)

type websocketWorker struct {
	client              v1.WebsocketHandlerClient
	registrationRequest *v1.RegistrationRequest
	handler             handlers.Handler[Ctx]
}
type websocketWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Handler             handlers.Handler[Ctx]
}

// Start implements Worker.
func (w *websocketWorker) Start(ctx context.Context) error {
	initReq := &v1.ClientMessage{
		Content: &v1.ClientMessage_RegistrationRequest{
			RegistrationRequest: w.registrationRequest,
		},
	}

	createStream := func(ctx context.Context) (workers.Stream[v1.ClientMessage, v1.RegistrationResponse, *v1.ServerMessage], error) {
		return w.client.HandleEvents(ctx)
	}

	handlerSrvMsg := func(msg *v1.ServerMessage) (*v1.ClientMessage, error) {
		if msg.GetWebsocketEventRequest() != nil {
			handlerCtx := NewCtx(msg)

			err := w.handler(handlerCtx)
			if err != nil {
				handlerCtx.WithError(err)
			}
			return handlerCtx.ToClientMessage(), nil
		}

		return nil, errors.NewWithCause(
			codes.Internal,
			"WebsocketWorker: Unhandled server message",
			errorsstd.New("unhandled server message"),
		)
	}

	return workers.HandleStream(ctx, createStream, initReq, handlerSrvMsg)
}

func newWebsocketWorker(opts *websocketWorkerOpts) *websocketWorker {
	conn, err := grpcx.GetConnection()
	if err != nil {
		panic(errors.NewWithCause(
			codes.Unavailable,
			"NewWebsocketWorker: Unable to reach WebsocketHandlerClient",
			err,
		))
	}

	client := v1.NewWebsocketHandlerClient(conn)

	return &websocketWorker{
		client:              client,
		registrationRequest: opts.RegistrationRequest,
		handler:             opts.Handler,
	}
}
