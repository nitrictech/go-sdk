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

package topics

import (
	"context"

	errorsstd "errors"

	grpcx "github.com/nitrictech/go-sdk/internal/grpc"
	"github.com/nitrictech/go-sdk/internal/handlers"
	"github.com/nitrictech/go-sdk/nitric/errors"
	"github.com/nitrictech/go-sdk/nitric/errors/codes"
	"github.com/nitrictech/go-sdk/nitric/workers"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

type subscriptionWorker struct {
	client              v1.SubscriberClient
	registrationRequest *v1.RegistrationRequest
	handler             handlers.Handler[Ctx]
}
type subscriptionWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Handler             handlers.Handler[Ctx]
}

// Start implements Worker.
func (s *subscriptionWorker) Start(ctx context.Context) error {
	initReq := &v1.ClientMessage{
		Content: &v1.ClientMessage_RegistrationRequest{
			RegistrationRequest: s.registrationRequest,
		},
	}

	createStream := func(ctx context.Context) (workers.Stream[v1.ClientMessage, v1.RegistrationResponse, *v1.ServerMessage], error) {
		return s.client.Subscribe(ctx)
	}

	handleSrvMsg := func(msg *v1.ServerMessage) (*v1.ClientMessage, error) {
		if msg.GetMessageRequest() != nil {
			handlerCtx := NewCtx(msg)

			err := s.handler(handlerCtx)
			if err != nil {
				handlerCtx.WithError(err)
			}

			return handlerCtx.ToClientMessage(), nil
		}

		return nil, errors.NewWithCause(
			codes.Internal,
			"SubscriptionWorker: Unhandled server message",
			errorsstd.New("unhandled server message"),
		)
	}

	return workers.HandleStream(ctx, createStream, initReq, handleSrvMsg)
}

func newSubscriptionWorker(opts *subscriptionWorkerOpts) *subscriptionWorker {
	conn, err := grpcx.GetConnection()
	if err != nil {
		panic(errors.NewWithCause(
			codes.Unavailable,
			"NewSubscriptionWorker: Unable to reach SubscriberClient",
			err,
		))
	}

	client := v1.NewSubscriberClient(conn)

	return &subscriptionWorker{
		client:              client,
		registrationRequest: opts.RegistrationRequest,
		handler:             opts.Handler,
	}
}
