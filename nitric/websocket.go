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
	registerResult := <-defaultManager.registerResource(&resourcesv1.ResourceDeclareRequest{
		Id: &resourcesv1.ResourceIdentifier{
			Type: resourcesv1.ResourceType_Websocket,
			Name: name,
		},
	})
	if registerResult.Err != nil {
		return nil, registerResult.Err
	}

	actions := []resourcesv1.Action{resourcesv1.Action_WebsocketManage}

	m, err := defaultManager.registerPolicy(registerResult.Identifier, actions...)
	if err != nil {
		return nil, err
	}

	wClient := websocketsv1.NewWebsocketClient(m.conn)

	return &websocket{
		manager: defaultManager,
		client:  wClient,
		name:    name,
	}, nil
}

func (w *websocket) Name() string {
	return w.name
}

func (w *websocket) On(eventType handler.WebsocketEventType, middleware ...handler.WebsocketMiddleware) {
	// mapping handler.WebsocketEventType to protobuf requirement i.e websocketsv1.WebsocketEventType
	var _eventType websocketsv1.WebsocketEventType
	switch eventType {
	case handler.WebsocketDisconnect:
		_eventType = websocketsv1.WebsocketEventType_Disconnect
	case handler.WebsocketMessage:
		_eventType = websocketsv1.WebsocketEventType_Message
	default:
		_eventType = websocketsv1.WebsocketEventType_Connect
	}

	registrationRequest := &websocketsv1.RegistrationRequest{
		SocketName: w.name,
		EventType:  _eventType,
	}
	composeHandler := handler.ComposeWebsocketMiddleware(middleware...)

	opts := &workers.WebsocketWorkerOpts{
		RegistrationRequest: registrationRequest,
		Middleware:          composeHandler,
	}

	worker := workers.NewWebsocketWorker(opts)
	w.manager.addWorker("WebsocketWorker:"+strings.Join([]string{
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
