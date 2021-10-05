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

type DataRequest interface {
	Data() []byte
	MimeType() string
}

type dataRequestImpl struct {
	data     []byte
	mimeType string
}

func (d *dataRequestImpl) Data() []byte {
	return d.data
}

func (d *dataRequestImpl) MimeType() string {
	return d.mimeType
}

type HttpRequest interface {
	DataRequest
	Method() string
	Path() string
	Query() map[string][]string
	Headers() map[string][]string
}

type httpRequestImpl struct {
	dataRequestImpl
	method  string
	path    string
	query   map[string][]string
	headers map[string][]string
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
