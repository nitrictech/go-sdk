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

package websockets

import (
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
)

type EventType string

var EventTypes = []EventType{EventType_Connect, EventType_Disconnect, EventType_Message}

const (
	EventType_Connect    EventType = "connect"
	EventType_Disconnect EventType = "disconnect"
	EventType_Message    EventType = "message"
)

type Request interface {
	SocketName() string
	EventType() EventType
	ConnectionID() string
	QueryParams() map[string][]string
	Message() string
}

type requestImpl struct {
	socketName   string
	eventType    EventType
	connectionId string
	queryParams  map[string]*websocketspb.QueryValue
	message      string
}

func (w *requestImpl) SocketName() string {
	return w.socketName
}

func (w *requestImpl) EventType() EventType {
	return w.eventType
}

func (w *requestImpl) ConnectionID() string {
	return w.connectionId
}

func (w *requestImpl) QueryParams() map[string][]string {
	queryParams := map[string][]string{}

	for k, v := range w.queryParams {
		queryParams[k] = v.Value
	}

	return queryParams
}

func (w *requestImpl) Message() string {
	return w.message
}
