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

package workers

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/handler"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
)

type WebsocketWorker struct {
	client              v1.WebsocketHandlerClient
	registrationRequest *v1.RegistrationRequest
	middleware          handler.WebsocketMiddleware
}
type WebsocketWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Middleware          handler.WebsocketMiddleware
}

// Start implements Worker.
func (w *WebsocketWorker) Start(ctx context.Context) error {
	initReq := &v1.ClientMessage{
		Content: &v1.ClientMessage_RegistrationRequest{
			RegistrationRequest: w.registrationRequest,
		},
	}

	// Create the request stream and send the initial request
	stream, err := w.client.HandleEvents(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(initReq)
	if err != nil {
		return err
	}
	for {
		var ctx *handler.WebsocketContext

		resp, err := stream.Recv()

		if err == io.EOF {
			err = stream.CloseSend()
			if err != nil {
				return err
			}

			return nil
		} else if err == nil && resp.GetRegistrationResponse() != nil {
			// Blob Notification has connected with Nitric server
			fmt.Println("WebsocketWorker connected with Nitric server")
		} else if err == nil && resp.GetWebsocketEventRequest() != nil {
			ctx = handler.NewWebsocketContext(resp)
			ctx, err = w.middleware(ctx, handler.WebsocketDummy)
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

func NewWebsocketWorker(opts *WebsocketWorkerOpts) *WebsocketWorker {
	ctx, _ := context.WithTimeout(context.TODO(), constants.NitricDialTimeout())

	conn, err := grpc.DialContext(
		ctx,
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		panic(errors.NewWithCause(
			codes.Unavailable,
			"NewWebsocketWorker: Unable to reach StorageListenerClient",
			err,
		))
	}

	client := v1.NewWebsocketHandlerClient(conn)

	return &WebsocketWorker{
		client:              client,
		registrationRequest: opts.RegistrationRequest,
		middleware:          opts.Middleware,
	}
}
