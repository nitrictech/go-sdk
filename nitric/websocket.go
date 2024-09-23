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

	"github.com/nitrictech/go-sdk/nitric/websockets"
	resourcesv1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	websocketsv1 "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
)

// Websocket - Nitric Websocket API Resource
type Websocket interface {
	// Name - Get the name of the Websocket API
	Name() string
	// On registers a handler for a specific event type on the websocket
	// Valid function signatures for middleware are:
	//
	//	func()
	//	func() error
	//	func(*websocket.Ctx)
	//	func(*websocket.Ctx) error
	//	func(*websocket.Ctx) *websocket.Ctx
	//	func(*websocket.Ctx) (*websocket.Ctx, error)
	//	func(*websocket.Ctx, Handler[websocket.Ctx]) *websocket.Ctx
	//	func(*websocket.Ctx, Handler[websocket.Ctx]) error
	//	func(*websocket.Ctx, Handler[websocket.Ctx]) (*websocket.Ctx, error)
	//	Middleware[websocket.Ctx]
	//	Handler[websocket.Ctx]
	On(eventType websockets.EventType, mwares ...interface{})
	// Send a message to a specific connection
	Send(ctx context.Context, connectionId string, message []byte) error
	// Close a specific connection
	Close(ctx context.Context, connectionId string) error
}

type websocket struct {
	Websocket

	name    string
	manager *manager
	client  websocketsv1.WebsocketClient
}

// NewWebsocket - Create a new Websocket API resource
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

	err := defaultManager.registerPolicy(registerResult.Identifier, actions...)
	if err != nil {
		return nil, err
	}

	wClient := websocketsv1.NewWebsocketClient(defaultManager.conn)

	return &websocket{
		manager: defaultManager,
		client:  wClient,
		name:    name,
	}, nil
}

func (w *websocket) Name() string {
	return w.name
}

func (w *websocket) On(eventType websockets.EventType, middleware ...interface{}) {
	var _eventType websocketsv1.WebsocketEventType
	switch eventType {
	case websockets.EventType_Disconnect:
		_eventType = websocketsv1.WebsocketEventType_Disconnect
	case websockets.EventType_Message:
		_eventType = websocketsv1.WebsocketEventType_Message
	default:
		_eventType = websocketsv1.WebsocketEventType_Connect
	}

	registrationRequest := &websocketsv1.RegistrationRequest{
		SocketName: w.name,
		EventType:  _eventType,
	}

	middlewares, err := interfacesToMiddleware[websockets.Ctx](middleware)
	if err != nil {
		panic(err)
	}

	composeHandler := ComposeMiddleware(middlewares...)

	opts := &websocketWorkerOpts{
		RegistrationRequest: registrationRequest,
		Middleware:          composeHandler,
	}

	worker := newWebsocketWorker(opts)
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
