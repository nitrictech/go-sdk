package job

import (
	batch "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
)

type Context struct {
	id       string
	Request  Request
	Response *Response
	Extras   map[string]interface{}
	Next     Middleware
}

func (c *Context) ToClientMessage() *batch.ClientMessage {
	return &batch.ClientMessage{
		Id: c.id,
		Content: &batch.ClientMessage_JobResponse{
			JobResponse: &batch.JobResponse{
				Success: true,
			},
		},
	}
}

func NewJobContext(msg *batch.ServerMessage) *Context {
	return &Context{
		id: msg.Id,
		Request: &request{
			jobName: msg.GetJobRequest().JobName,
			data:    msg.GetJobRequest().Data.GetStruct().AsMap(),
		},
		Response: &Response{
			Success: true,
		},
	}
}

func (c *Context) WithError(err error) {
	c.Response = &Response{
		Success: false,
	}
}
