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

type HandlerBuilder interface {
	Http(...HttpMiddleware) HandlerBuilder
	Event(...EventMiddleware) HandlerBuilder
	Default(...TriggerMiddleware) HandlerBuilder
	Start() error
}

type HandlerProvider interface {
	GetHttp() HttpMiddleware
	GetEvent() EventMiddleware
	GetDefault() TriggerMiddleware
}

type faasClientImpl struct {
	http  HttpMiddleware
	event EventMiddleware
	trig  TriggerMiddleware
}

func (f *faasClientImpl) Http(mwares ...HttpMiddleware) HandlerBuilder {
	f.http = ComposeHttpMiddlware(mwares...)
	return f
}

func (f *faasClientImpl) GetHttp() HttpMiddleware {
	return f.http
}

func (f *faasClientImpl) Event(mwares ...EventMiddleware) HandlerBuilder {
	f.event = ComposeEventMiddleware(mwares...)
	return f
}

func (f *faasClientImpl) GetEvent() EventMiddleware {
	return f.event
}

func (f *faasClientImpl) Default(mwares ...TriggerMiddleware) HandlerBuilder {
	f.trig = ComposeTriggerMiddleware(mwares...)
	return f
}

func (f *faasClientImpl) GetDefault() TriggerMiddleware {
	return f.trig
}

func (f *faasClientImpl) Start() error {
	// Fail if no handlers were provided
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

	fsc := pb.NewFaasServiceClient(conn)

	return f.startWithClient(fsc)
}

func (f *faasClientImpl) startWithClient(fsc pb.FaasServiceClient) error {
	if f.http == nil && f.event == nil && f.trig == nil {
		return fmt.Errorf("no valid handlers provided")
	}

	if stream, err := fsc.TriggerStream(context.TODO()); err == nil {
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

// Creates a new HandlerBuilder
func New() HandlerBuilder {
	return &faasClientImpl{}
}
