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

	pb "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
)

func faasLoop(stream pb.FaasService_TriggerStreamClient, f HandlerProvider, errorCh chan error) {
	for {
		// Block receiving a message
		srvrMsg, err := stream.Recv()

		if err != nil {
			// TODO: Make sure we use the correct kind of error types here
			errorCh <- errors.FromGrpcError(err)
			break
		}

		clientMsg := &pb.ClientMessage{
			Id: srvrMsg.GetId(),
		}

		// We have a trigger
		if srvrMsg.GetTriggerRequest() != nil {
			ctx, err := triggerContextFromGrpcTriggerRequest(srvrMsg.GetTriggerRequest())

			if err != nil {
				fmt.Println("There was an error reading the TriggerRequest", err)
				// Return a bad request here...
				errorCh <- errors.NewWithCause(
					codes.Internal,
					"faasLoop: error reading the TriggerRequest",
					err,
				)
				break
			}

			// Interrogate the HandlerProvider to see if it is capable of handling this trigger
			var funcErr error = nil
			if ctx.Http() != nil && f.GetHttp() != nil {
				// handle http
				ctx.http, funcErr = f.GetHttp()(ctx.Http(), httpDummy)

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
				if ctx.http != nil {
					ctx.http.Response.Body = []byte("Internal Server Error")
					ctx.http.Response.Headers = map[string][]string{
						"Content-Type": {"text/plain"},
					}
					ctx.http.Response.Status = 500
				} else if ctx.event != nil {
					ctx.event.Response.Success = false
				}
			}

			triggerResponse, err := triggerContextToGrpcTriggerResponse(ctx)

			if err != nil {
				fmt.Println("Error translating handler response", err)
				errorCh <- errors.FromGrpcError(err)
				break
			}

			clientMsg.Content = &pb.ClientMessage_TriggerResponse{
				TriggerResponse: triggerResponse,
			}

			if err := stream.Send(clientMsg); err != nil {
				fmt.Println("Failed to send msg", err)
				errorCh <- errors.FromGrpcError(err)
				break
			}
		} else if srvrMsg.GetInitResponse() != nil {
			fmt.Println("Function connected to membrane", err)
		}
	}
}
