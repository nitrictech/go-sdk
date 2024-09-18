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

package nitric

import (
	"context"
	errorsstd "errors"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	httpx "github.com/nitrictech/go-sdk/api/http"
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
)

type apiWorker struct {
	client              v1.ApiClient
	middleware          Middleware[httpx.Ctx]
	registrationRequest *v1.RegistrationRequest
}

type apiWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Middleware          Middleware[httpx.Ctx]
}

var _ streamWorker = (*apiWorker)(nil)

// Start implements Worker.
func (a *apiWorker) Start(ctx context.Context) error {
	initReq := &v1.ClientMessage{
		Content: &v1.ClientMessage_RegistrationRequest{
			RegistrationRequest: a.registrationRequest,
		},
	}

	stream, err := a.client.Serve(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(initReq)
	if err != nil {
		return err
	}

	for {
		var ctx *httpx.Ctx

		resp, err := stream.Recv()

		if errorsstd.Is(err, io.EOF) {
			err = stream.CloseSend()
			if err != nil {
				return err
			}

			return nil
		} else if err == nil && resp.GetRegistrationResponse() != nil {
			// Function connected with Nitric server
			fmt.Println("Function connected with Nitric server")
		} else if err == nil && resp.GetHttpRequest() != nil {
			ctx = httpx.NewCtx(resp)

			ctx, err = a.middleware(ctx, dummyHandler)
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

func newApiWorker(opts *apiWorkerOpts) *apiWorker {
	conn, err := grpc.NewClient(constants.NitricAddress(), constants.DefaultOptions()...)
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
		middleware:          opts.Middleware,
	}
}
