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

type FileCtx struct {
	Request  FileRequest
	Response *FileResponse
	Extras   map[string]interface{}
}
