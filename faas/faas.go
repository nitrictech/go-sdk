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
	"io"

	"github.com/nitrictech/go-sdk/constants"
	pb "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

// NitricFunction - a function built using Nitric, to be executed
type NitricFunction func(*NitricTrigger) (*NitricResponse, error)

func faasLoop(stream pb.Faas_TriggerStreamClient, f NitricFunction, errorCh chan error) {
	for {
		// Block recieving a message
		srvrMsg, err := stream.Recv()
		clientMsg := &pb.ClientMessage{
			Id: srvrMsg.GetId(),
		}

		if err != nil {
			// TODO: Make sure we use the correct kind of error types here
			errorCh <- err
			break
		}

		// We have a trigger
		if srvrMsg.GetTriggerRequest() != nil {
			req, err := FromGrpcTriggerRequest(srvrMsg.GetTriggerRequest())

			if err != nil {
				fmt.Println("There was an error reading the TriggerRequest", err)
				// Return a bad request here...

				continue
			}
			// Let the membrane know the function is ready for initializatio
			// Process this trigger
			response, err := f(req)

			if err != nil {
				fmt.Println("Function return an error", err)
				// Return an error here...
				response = req.DefaultResponse()
				response.SetData([]byte("Internal Error"))

				if response.context.IsHttp() {
					http := response.context.AsHttp()
					http.Headers = map[string]string{
						"Content-Type": "text/plain",
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
				if err != io.EOF {
					fmt.Println("Failed to send msg", err)
					errorCh <- err
					break
				}
				fmt.Println("EOF encountered from server", err)
				continue
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
	conn, err := grpc.Dial(constants.NitricAddress(), grpc.WithInsecure())

	if err != nil {
		return err
	}

	faasClient := pb.NewFaasClient(conn)

	return StartWithClient(f, faasClient)
}

func StartWithClient(f NitricFunction, faasClient pb.FaasClient) error {
	if stream, err := faasClient.TriggerStream(context.TODO()); err == nil {
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
