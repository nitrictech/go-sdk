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

package context

type HttpContext struct {
	Request  HttpRequest
	Response *HttpResponse
	Extras   map[string]interface{}
}

type MessageContext struct {
	Request  MessageRequest
	Response *MessageResponse
	Extras   map[string]interface{}
}

type IntervalContext struct {
	Request  IntervalRequest
	Response *IntervalResponse
	Extras   map[string]interface{}
}

type BlobEventContext struct {
	Request  BlobEventRequest
	Response *BlobEventResponse
	Extras   map[string]interface{}
}

type FileEventContext struct {
	Request  FileEventRequest
	Response *FileEventResponse
	Extras   map[string]interface{}
}

type WebsocketContext struct {
	Request  WebsocketRequest
	Response *WebsocketResponse
	Extras   map[string]interface{}
}
