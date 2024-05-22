package workers

import (
	"context"
	"fmt"
	"io"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/handler"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	"google.golang.org/grpc"
)

type ApiWorker struct {
	client              v1.ApiClient
	middleware          handler.HttpMiddleware
	registrationRequest *v1.RegistrationRequest
}

type ApiWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Middleware          handler.HttpMiddleware
}

var _ Worker = (*ApiWorker)(nil)

// Start implements Worker.
func (a *ApiWorker) Start(ctx context.Context) error {
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
		var ctx *handler.HttpContext

		resp, err := stream.Recv()

		if err == io.EOF {
			err = stream.CloseSend()
			if err != nil {
				return err
			}

			return nil
		} else if err == nil && resp.GetRegistrationResponse() != nil {
			// Function connected with Nitric server
			fmt.Println("Function connected with Nitric server")
		} else if err == nil && resp.GetHttpRequest() != nil {
			ctx = handler.NewHttpContext(resp)

			ctx, err = a.middleware(ctx, handler.HttpDummy)
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

func NewApiWorker(opts *ApiWorkerOpts) *ApiWorker {
	ctx, _ := context.WithTimeout(context.TODO(), constants.NitricDialTimeout())

	conn, err := grpc.DialContext(
		ctx,
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		panic(errors.NewWithCause(
			codes.Unavailable,
			"NewApiWorker: Unable to reach ApiClient",
			err,
		))
	}

	client := v1.NewApiClient(conn)

	return &ApiWorker{
		client:     client,
		registrationRequest: opts.RegistrationRequest,
		middleware: opts.Middleware,
	}
}
