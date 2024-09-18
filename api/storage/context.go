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

package storage

import storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"

type Ctx struct {
	id       string
	Request  Request
	Response *Response
	Extras   map[string]interface{}
}

func (c *Ctx) ToClientMessage() *storagepb.ClientMessage {
	return &storagepb.ClientMessage{
		Id: c.id,
		Content: &storagepb.ClientMessage_BlobEventResponse{
			BlobEventResponse: &storagepb.BlobEventResponse{
				Success: c.Response.Success,
			},
		},
	}
}

func NewCtx(msg *storagepb.ServerMessage) *Ctx {
	req := msg.GetBlobEventRequest()

	return &Ctx{
		id: msg.Id,
		Request: &requestImpl{
			key: req.GetBlobEvent().Key,
		},
		Response: &Response{
			Success: true,
		},
	}
}

func (c *Ctx) WithError(err error) {
	c.Response = &Response{
		Success: false,
	}
}

// type FileCtx struct {
// 	Request  FileRequest
// 	Response *FileResponse
// 	Extras   map[string]interface{}
// }
