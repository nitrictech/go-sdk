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

package http

type Request interface {
	Method() string
	Path() string
	Data() []byte
	Query() map[string][]string
	Headers() map[string][]string
	PathParams() map[string]string
}

type RequestImpl struct {
	method     string
	path       string
	data       []byte
	query      map[string][]string
	headers    map[string][]string
	pathParams map[string]string
}

func (h *RequestImpl) Method() string {
	return h.method
}

func (h *RequestImpl) Path() string {
	return h.path
}

func (h *RequestImpl) Data() []byte {
	return h.data
}

func (h *RequestImpl) Query() map[string][]string {
	return h.query
}

func (h *RequestImpl) Headers() map[string][]string {
	return h.headers
}

func (h *RequestImpl) PathParams() map[string]string {
	return h.pathParams
}