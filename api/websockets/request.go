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
