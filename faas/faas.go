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
	"context"
	"fmt"

	pb "github.com/nitrictech/apis/go/nitric/v1"
	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	"google.golang.org/grpc"
)

// NitricFunction - a function built using Nitric, to be executed
type NitricFunction func(*NitricTrigger) (*NitricResponse, error)

func faasLoop(stream pb.FaasService_TriggerStreamClient, f NitricFunction, errorCh chan error) {

	for {
		// Block receiving a message
		srvrMsg, err := stream.Recv()
		clientMsg := &pb.ClientMessage{
			Id: srvrMsg.GetId(),
		}

		if err != nil {
			// TODO: Make sure we use the correct kind of error types here
			errorCh <- errors.FromGrpcError(err)
			break
		}

		// We have a trigger
		if srvrMsg.GetTriggerRequest() != nil {
			req, err := FromGrpcTriggerRequest(srvrMsg.GetTriggerRequest())

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
			// Let the membrane know the function is ready for initialization
			// Process this trigger
			response, err := f(req)

			if err != nil {
				fmt.Println("Function returned an error", err)
				// Return an error here...
				response = req.DefaultResponse()
				response.SetData([]byte("Internal Error"))

				if response.context.IsHttp() {
					http := response.context.AsHttp()
					http.Headers = map[string][]string{
						"Content-Type": {"text/plain"},
					}
					http.Status = 500
					// internal server error
				} else if response.context.IsTopic() {
					topic := response.context.AsTopic()
					topic.Success = false
					// mark as unsuccessful here...
				}
			}

			triggerResponse := response.ToTriggerResponse()

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

// Start - Starts accepting requests for the provided NitricFunction
// Begins streaming using the default Nitric FaaS gRPC client
// This should be the only method called in the 'main' method of your entrypoint package
func Start(f NitricFunction) error {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)

	if err != nil {
		return errors.NewWithCause(
			codes.Unavailable,
			"faas.Start: Unable to reach FaasServiceServer",
			err,
		)
	}

	FaasServiceClient := pb.NewFaasServiceClient(conn)

	return StartWithClient(f, FaasServiceClient)
}

func StartWithClient(f NitricFunction, FaasServiceClient pb.FaasServiceClient) error {
	if stream, err := FaasServiceClient.TriggerStream(context.TODO()); err == nil {
		// Let the membrane know the function is ready for initialization
		err := stream.Send(&pb.ClientMessage{
			Content: &pb.ClientMessage_InitRequest{
				InitRequest: &pb.InitRequest{},
			},
		})

		if err != nil {
			return err
		}

		errChan := make(chan error)

		// Start faasLoop in a go routine
		go faasLoop(stream, f, errChan)

		return <-errChan
	} else {
		return err
	}
}
