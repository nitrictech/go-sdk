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

package resources

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/faas"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	websocketv1 "github.com/nitrictech/nitric/core/pkg/api/nitric/websocket/v1"
)

type Websocket interface {
	Name() string
	On(eventType faas.WebsocketEventType, mwares ...faas.WebsocketMiddleware)
	Send(ctx context.Context, connectionId string, message []byte) error
	Close(ctx context.Context, connectionId string) error
	URL(ctx context.Context) (string, error)
}

type websocket struct {
	Websocket

	name    string
	manager Manager
	client  websocketv1.WebsocketServiceClient
}

// NewCollection register this collection as a required resource for the calling function/container.
func NewWebsocket(name string) (Websocket, error) {
	return defaultManager.newWebsocket(name)
}

func (m *manager) newWebsocket(name string) (Websocket, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)
	if err != nil {
		return nil, errors.NewWithCause(
			codes.Unavailable,
			"Websocket.New: Unable to reach WebsocketServiceServer",
			err,
		)
	}

	wClient := websocketv1.NewWebsocketServiceClient(conn)

	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	res := &v1.Resource{
		Type: v1.ResourceType_Websocket,
		Name: name,
	}

	dr := &v1.ResourceDeclareRequest{
		Resource: res,
	}

	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []v1.Action{v1.Action_WebsocketManage}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(res, actions))
	if err != nil {
		return nil, err
	}

	return &websocket{
		manager: m,
		client:  wClient,
		name:    name,
	}, nil
}

func (w *websocket) Name() string {
	return w.name
}

func (w *websocket) On(eventType faas.WebsocketEventType, middleware ...faas.WebsocketMiddleware) {
	f := faas.New()

	f.Websocket(middleware...).
		WithWebsocketWorkerOptions(faas.WebsocketWorkerOptions{
			Socket:    w.name,
			EventType: eventType,
		})

	w.manager.addWorker(fmt.Sprintf("websocket:%s/%s", w.name, eventType), f)
}

func (w *websocket) Send(ctx context.Context, connectionId string, message []byte) error {
	_, err := w.client.Send(ctx, &websocketv1.WebsocketSendRequest{
		Socket:       w.name,
		ConnectionId: connectionId,
		Data:         message,
	})

	return err
}

func (w *websocket) Close(ctx context.Context, connectionId string) error {
	_, err := w.client.Close(ctx, &websocketv1.WebsocketCloseRequest{
		Socket:       w.name,
		ConnectionId: connectionId,
	})

	return err
}

func (w *websocket) details(ctx context.Context) (*v1.ResourceDetailsResponse, error) {
	rsc, err := w.manager.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	dec := &v1.ResourceDetailsRequest{
		Resource: &v1.Resource{
			Name: w.name,
			Type: v1.ResourceType_Websocket,
		},
	}

	return rsc.Details(ctx, dec)
}

func (w *websocket) URL(ctx context.Context) (string, error) {
	resp, err := w.details(ctx)
	if err != nil {
		return "", err
	}

	return resp.GetWebsocket().GetUrl(), nil
}
