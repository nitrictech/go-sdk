package workers

import (
	"context"
	"fmt"
	"io"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/handler"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	"google.golang.org/grpc"
)


type WebsocketWorker struct {
	client v1.WebsocketHandlerClient
	registrationRequest *v1.RegistrationRequest 
	middleware handler.WebsocketMiddleware 
}


type WebsocketWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Middleware handler.WebsocketMiddleware
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
	if err != nil{
		return err
	}

	err = stream.Send(initReq)
	if err != nil{
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
		client: client,
		registrationRequest: opts.RegistrationRequest,
		middleware: opts.Middleware,
	}
}