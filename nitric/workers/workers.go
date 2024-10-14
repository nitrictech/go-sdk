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
	"errors"
	"fmt"
	"io"

	"google.golang.org/grpc"
)

type StreamWorker interface {
	Start(context.Context) error
}

type StdServerMsg[RegistrationResponse any] interface {
	GetRegistrationResponse() *RegistrationResponse
}

type Stream[ClientMessage any, RegistrationResponse any, ServerMessage StdServerMsg[RegistrationResponse]] interface {
	Send(*ClientMessage) error
	Recv() (ServerMessage, error)
	grpc.ClientStream
}

// HandleStream runs a nitric worker, in the standard request/response pattern.
// No changes needed here other than the updated types in the signature.
func HandleStream[ClientMessage any, RegistrationResponse any, ServerMessage StdServerMsg[RegistrationResponse]](
	ctx context.Context,
	createStream func(ctx context.Context) (Stream[ClientMessage, RegistrationResponse, ServerMessage], error),
	initReq *ClientMessage,
	handleServerMsg func(msg ServerMessage) (*ClientMessage, error),
) error {
	stream, err := createStream(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(initReq)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Context canceled, closing stream\n")
			// If the context is canceled, close the stream and return
			err := stream.CloseSend()
			if err != nil {
				return err
			}
			return nil

		default:
			// Receive the next message
			serverMsg, err := stream.Recv()

			if errors.Is(err, io.EOF) {
				// Close the stream and exit normally on EOF
				err = stream.CloseSend()
				if err != nil {
					return err
				}
				return nil
			} else if err != nil {
				return err
			}

			if serverMsg.GetRegistrationResponse() != nil {
				// No need to respond to the registration responses (they're just acks)
				continue
			}

			clientMsg, err := handleServerMsg(serverMsg)
			if err != nil {
				return err
			}

			err = stream.Send(clientMsg)
			if err != nil {
				return err
			}
		}
	}
}
