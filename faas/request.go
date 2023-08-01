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

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type DataRequest interface {
	Data() []byte
	MimeType() string
	Context() context.Context
}

type dataRequestImpl struct {
	data         []byte
	mimeType     string
	traceContext map[string]string
}

func (d *dataRequestImpl) Data() []byte {
	return d.data
}

func (d *dataRequestImpl) MimeType() string {
	return d.mimeType
}

func (d *dataRequestImpl) Context() context.Context {
	phc := propagation.HeaderCarrier{}

	for k, v := range d.traceContext {
		phc.Set(k, v)
	}

	return otel.GetTextMapPropagator().Extract(context.Background(), phc)
}

type HttpRequest interface {
	DataRequest
	Context() context.Context
	Method() string
	Path() string
	Query() map[string][]string
	Headers() map[string][]string
	PathParams() map[string]string
}

type httpRequestImpl struct {
	dataRequestImpl
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

type EventRequest interface {
	DataRequest
	Topic() string
}

type eventRequestImpl struct {
	dataRequestImpl
	topic string
}

func (e *eventRequestImpl) Topic() string {
	return e.topic
}

type BucketNotificationRequest interface {
	Key() string
	NotificationType() NotificationType
}

type bucketNotificationRequestImpl struct {
	key              string
	notificationType NotificationType
}

func (b *bucketNotificationRequestImpl) Key() string {
	return b.key
}

func (b *bucketNotificationRequestImpl) NotificationType() NotificationType {
	return b.notificationType
}

type WebsocketRequest interface {
	DataRequest

	Socket() string
	EventType() WebsocketEventType
	ConnectionID() string
	QueryParams() map[string][]string
}

type websocketRequestImpl struct {
	dataRequestImpl

	socket string
	eventType WebsocketEventType
	connectionId string
	queryParams map[string]*v1.QueryValue
}

func (w *websocketRequestImpl) Socket() string {
	return w.socket
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
