// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package faas

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	pb "github.com/nitrictech/apis/go/nitric/v1"
	apierrors "github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
)

const concurrentRequestLimit int = 10

type processOneArgs struct {
	ctx      *triggerContextImpl
	stream   pb.FaasService_TriggerStreamClient
	f        HandlerProvider
	svrMsgID string
}

func faasLoop(stream pb.FaasService_TriggerStreamClient, f HandlerProvider, errorCh chan error) {
	wg := NewWorkPool(concurrentRequestLimit)

	for {
		if wg.Err() != nil {
			break
		}

		// Block receiving a message
		srvrMsg, err := stream.Recv()
		if err != nil {
			// TODO: Make sure we use the correct kind of error types here
			if errors.Is(err, io.EOF) {
				wg.AddError(err)
			} else {
				wg.AddError(apierrors.FromGrpcError(err))
			}
			break
		}

		// We have a trigger
		if srvrMsg.GetTriggerRequest() != nil {
			ctx, err := triggerContextFromGrpcTriggerRequest(srvrMsg.GetTriggerRequest())
			if err != nil {
				fmt.Println("There was an error reading the TriggerRequest", err)
				// Return a bad request here...
				wg.AddError(apierrors.NewWithCause(
					codes.Internal,
					"faasLoop: error reading the TriggerRequest",
					err,
				))
				break
			}

			wg.Go(func(a interface{}) error {
				aa := a.(*processOneArgs)

				return faasProcessOne(aa.ctx, aa.stream, aa.f, aa.svrMsgID)
			}, &processOneArgs{ctx: ctx,
				stream:   stream,
				f:        f,
				svrMsgID: srvrMsg.GetId(),
			})
		} else if srvrMsg.GetInitResponse() != nil {
			fmt.Println("Function connected to membrane")
		}
	}

	wg.Wait()

	errorCh <- wg.Err()
}

func withInternalServerError(ctx *triggerContextImpl) *triggerContextImpl {
	if ctx.http != nil {
		ctx.http.Response.Body = []byte("Internal Server Error")
		ctx.http.Response.Headers = map[string][]string{
			"Content-Type": {"text/plain"},
		}
		ctx.http.Response.Status = 500
	} else if ctx.event != nil {
		ctx.event.Response.Success = false
	}

	return ctx
}

func faasPanicRecovery(ctx *triggerContextImpl, stream pb.FaasService_TriggerStreamClient, f HandlerProvider, svrMsgID string) {
	if rErr := recover(); rErr != nil {
		fmt.Println(errors.WithStack(fmt.Errorf("the handler function paniced: %v", rErr)))

		err := faasSendResponse(withInternalServerError(ctx), stream, svrMsgID)
		if err != nil {
			fmt.Println(errors.WithMessage(err, "Error sending error response"))
		}
	}
}

func faasProcessOne(ctx *triggerContextImpl, stream pb.FaasService_TriggerStreamClient, f HandlerProvider, svrMsgID string) (err error) {
	// try our best to send an error back to the client if a handler panics.
	defer faasPanicRecovery(ctx, stream, f, svrMsgID)

	// Interrogate the HandlerProvider to see if it is capable of handling this trigger
	var funcErr error = nil
	if ctx.Http() != nil && f.GetHttp(ctx.http.Request.Method()) != nil {
		// handle http
		ctx.http, funcErr = f.GetHttp(ctx.http.Request.Method())(ctx.Http(), httpDummy)

		if ctx.http == nil && funcErr == nil {
			funcErr = fmt.Errorf("nil context returned from http handler")
		}
	} else if ctx.Event() != nil && f.GetEvent() != nil {
		// handle event
		ctx.event, funcErr = f.GetEvent()(ctx.Event(), eventDummy)

		if ctx.event == nil && funcErr == nil {
			funcErr = fmt.Errorf("nil context returned from event handler")
		}
	} else if f.GetDefault() != nil {
		// handle trigger
		var newCtx TriggerContext
		newCtx, funcErr = f.GetDefault()(ctx, triggerDummy)
		if newCtx != nil {
			// Update response context
			ctx.event = newCtx.Event()
			ctx.http = newCtx.Http()
		} else if funcErr != nil {
			funcErr = fmt.Errorf("nil context returned from trigger handler")
		}
	} else {
		// No available handler for the event...
		// This is not a panic error case, but we need to return
		// a default unavailable response
		funcErr = fmt.Errorf("no handler available for trigger type")
	}

	if funcErr != nil {
		fmt.Printf("an error was returned by the handler function: %s\n", funcErr.Error())

		withInternalServerError(ctx)
	}

	return faasSendResponse(ctx, stream, svrMsgID)
}

func faasSendResponse(ctx *triggerContextImpl, stream pb.FaasService_TriggerStreamClient, svrMsgID string) error {
	triggerResponse, err := triggerContextToGrpcTriggerResponse(ctx)
	if err != nil {
		fmt.Println("Error translating handler response", err)
		return apierrors.FromGrpcError(err)
	}

	err = stream.Send(&pb.ClientMessage{
		Id: svrMsgID,
		Content: &pb.ClientMessage_TriggerResponse{
			TriggerResponse: triggerResponse,
		},
	})
	if err != nil {
		fmt.Println("Error sending handler response", err)
		return apierrors.FromGrpcError(err)
	}

	return nil
}
