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

package nitric

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/nitrictech/go-sdk/api/errors"
	"github.com/nitrictech/go-sdk/api/errors/codes"
	"github.com/nitrictech/go-sdk/constants"
	"github.com/nitrictech/go-sdk/handler"
	"github.com/nitrictech/go-sdk/workers"
	resourcesv1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	websocketsv1 "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
)

type Websocket interface {
	Name() string
	On(eventType handler.WebsocketEventType, mwares ...handler.WebsocketMiddleware)
	Send(ctx context.Context, connectionId string, message []byte) error
	Close(ctx context.Context, connectionId string) error
}

type websocket struct {
	Websocket

	name    string
	manager Manager
	client  websocketsv1.WebsocketClient
}

// NewCollection register this collection as a required resource for the calling function/container.
func NewWebsocket(name string) (Websocket, error) {
	return defaultManager.newWebsocket(name)
}

func (m *manager) newWebsocket(name string) (Websocket, error) {
	ctx, _ := context.WithTimeout(context.Background(), constants.NitricDialTimeout())

	conn, err := grpc.DialContext(
		ctx,
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

	rsc, err := m.resourceServiceClient()
	if err != nil {
		return nil, err
	}

	res := &resourcesv1.ResourceIdentifier{
		Type: resourcesv1.ResourceType_Websocket,
		Name: name,
	}

	dr := &resourcesv1.ResourceDeclareRequest{
		Id: res,
	}

	_, err = rsc.Declare(context.Background(), dr)
	if err != nil {
		return nil, err
	}

	actions := []resourcesv1.Action{resourcesv1.Action_WebsocketManage}

	_, err = rsc.Declare(context.Background(), functionResourceDeclareRequest(res, actions))
	if err != nil {
		return nil, err
	}

	wClient := websocketsv1.NewWebsocketClient(conn)

	return &websocket{
		manager: m,
		client:  wClient,
		name:    name,
	}, nil
}

func (w *websocket) Name() string {
	return w.name
}

func (w *websocket) On(eventType handler.WebsocketEventType, middleware ...handler.WebsocketMiddleware) {

	// mapping handler.WebsocketEventType to protobuf requirement i.e websocketsv1.WebsocketEventType
	var _eventType websocketsv1.WebsocketEventType;
	switch eventType {
	case  handler.WebsocketDisconnect:
		_eventType = websocketsv1.WebsocketEventType_Disconnect
	case handler.WebsocketMessage:
		_eventType = websocketsv1.WebsocketEventType_Message
	default:
		_eventType = websocketsv1.WebsocketEventType_Connect
	}

	registrationRequest := &websocketsv1.RegistrationRequest{
		SocketName: w.name,
		EventType: _eventType,
	}
	composeHandler := handler.ComposeWebsocketMiddleware(middleware...)

	opts := &workers.WebsocketWorkerOpts{
		RegistrationRequest: registrationRequest,
		Middleware: composeHandler,
	}

	worker := workers.NewWebsocketWorker(opts)
	w.manager.addWorker("WebsocketWorker:" + strings.Join([]string{
		w.name,
		string(eventType),
	}, "-"), worker)
}

func (w *websocket) Send(ctx context.Context, connectionId string, message []byte) error {
	_, err := w.client.SendMessage(ctx, &websocketsv1.WebsocketSendRequest{
		SocketName:   w.name,
		ConnectionId: connectionId,
		Data:         message,
	})

	return err
}

func (w *websocket) Close(ctx context.Context, connectionId string) error {
	_, err := w.client.CloseConnection(ctx, &websocketsv1.WebsocketCloseConnectionRequest{
		SocketName:   w.name,
		ConnectionId: connectionId,
	})

	return err
}
