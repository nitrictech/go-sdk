package workers

import (
	"context"
	"fmt"
	"io"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/handler"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	"google.golang.org/grpc"
)


type BlobEventWorker struct {
	client v1.StorageListenerClient
	registrationRequest *v1.RegistrationRequest 
	middleware handler.BlobEventMiddleware 
}


type BlobEventWorkerOpts struct {
	RegistrationRequest *v1.RegistrationRequest
	Middleware handler.BlobEventMiddleware
}

// Start implements Worker.
func (b *BlobEventWorker) Start(ctx context.Context) error {
	
	initReq := &v1.ClientMessage{
		Content: &v1.ClientMessage_RegistrationRequest{
			RegistrationRequest: b.registrationRequest,
		},
	}
	
	// Create the request stream and send the initial request
	stream, err := b.client.Listen(ctx)
	if err != nil{
		return err
	}

	err = stream.Send(initReq)
	if err != nil{
		return err
	}


	for {
		var ctx *handler.BlobEventContext

		resp, err := stream.Recv()

		if err == io.EOF {
			err = stream.CloseSend()
			if err != nil {
				return err
			}

			return nil
		} else if err == nil && resp.GetRegistrationResponse() != nil {
			// Blob Notification has connected with Nitric server
			fmt.Println("Function connected with Nitric server")
		} else if err == nil && resp.GetBlobEventRequest() != nil {
			
			ctx = handler.NewBlobEventContext(resp)
			ctx, err = b.middleware(ctx, handler.BlobEventDummy)

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

func NewBlobEventWorker(opts *BlobEventWorkerOpts) *BlobEventWorker {
	ctx, _ := context.WithTimeout(context.TODO(), constants.NitricDialTimeout())

	conn, err := grpc.DialContext(
		ctx,
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		panic(errors.NewWithCause(
			codes.Unavailable,
			"NewBlobEventWorker: Unable to reach StorageListenerClient",
			err,
		))
	}

	client := v1.NewStorageListenerClient(conn)
	
	return &BlobEventWorker{
		client: client,
		registrationRequest: opts.RegistrationRequest,
		middleware: opts.Middleware,
	}
}