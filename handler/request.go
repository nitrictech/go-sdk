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

package handler

import (
	"github.com/nitrictech/go-sdk/api/storage"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
)

// Http

type HttpRequest interface {
	Method() string
	Path() string
	Query() map[string][]string
	Headers() map[string][]string
	PathParams() map[string]string
}

type httpRequestImpl struct {
	method     string
	path       string
	query      map[string][]string
	headers    map[string][]string
	pathParams map[string]string
}

func (h *httpRequestImpl) Method() string {
	return h.method
}

func (h *httpRequestImpl) Path() string {
	return h.path
}

func (h *httpRequestImpl) Query() map[string][]string {
	return h.query
}

func (h *httpRequestImpl) Headers() map[string][]string {
	return h.headers
}

func (h *httpRequestImpl) PathParams() map[string]string {
	return h.pathParams
}

// Message

type MessageRequest interface {
	TopicName() string
	Message() map[string]interface{}
}

type messageRequestImpl struct {
	topicName string
	message   map[string]interface{}
}

func (m *messageRequestImpl) TopicName() string {
	return m.topicName
}

func (m *messageRequestImpl) Message() map[string]interface{} {
	return m.message
}

// Interval

type IntervalRequest interface {
	ScheduleName() string
}

type intervalRequestImpl struct {
	scheduleName string
}

func (i *intervalRequestImpl) ScheduleName() string {
	return i.scheduleName
}

// Blob Event

type BlobEventType string

var BlobEventTypes = []BlobEventType{WriteNotification, DeleteNotification}

const (
	WriteNotification  BlobEventType = "write"
	DeleteNotification BlobEventType = "delete"
)

type BlobEventRequest interface {
	Key() string
	NotificationType() BlobEventType
}

type blobEventRequestImpl struct {
	key              string
	notificationType BlobEventType
}

func (b *blobEventRequestImpl) Key() string {
	return b.key
}

func (b *blobEventRequestImpl) NotificationType() BlobEventType {
	return b.notificationType
}

// File Event

type FileEventRequest interface {
	Bucket() *storage.Bucket
	NotificationType() BlobEventType
}

type fileEventRequestImpl struct {
	bucket           storage.Bucket
	notificationType BlobEventType
}

func (f *fileEventRequestImpl) Bucket() storage.Bucket {
	return f.bucket
}

func (f *fileEventRequestImpl) NotificationType() BlobEventType {
	return f.notificationType
}

// Websocket

type WebsocketEventType string

var WebsocketEventTypes = []WebsocketEventType{WebsocketConnect, WebsocketDisconnect, WebsocketMessage}

const (
	WebsocketConnect    WebsocketEventType = "connect"
	WebsocketDisconnect WebsocketEventType = "disconnect"
	WebsocketMessage    WebsocketEventType = "message"
)

type WebsocketRequest interface {
	SocketName() string
	EventType() WebsocketEventType
	ConnectionID() string
	QueryParams() map[string][]string
	Message() string
}

type websocketRequestImpl struct {
	socketName   string
	eventType    WebsocketEventType
	connectionId string
	queryParams  map[string]*v1.QueryValue
	message      string
}

func (w *websocketRequestImpl) SocketName() string {
	return w.socketName
}

func (w *websocketRequestImpl) EventType() WebsocketEventType {
	return w.eventType
}

func (w *websocketRequestImpl) ConnectionID() string {
	return w.connectionId
}

func (w *websocketRequestImpl) QueryParams() map[string][]string {
	queryParams := map[string][]string{}

	for k, v := range w.queryParams {
		queryParams[k] = v.Value
	}

	return queryParams
}

func (w *websocketRequestImpl) Message() string {
	return w.message
}
