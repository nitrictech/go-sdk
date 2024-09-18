package topics

import topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"

type Ctx struct {
	id       string
	Request  Request
	Response *Response
	Extras   map[string]interface{}
}

func (c *Ctx) ToClientMessage() *topicspb.ClientMessage {
	return &topicspb.ClientMessage{
		Id: c.id,
		Content: &topicspb.ClientMessage_MessageResponse{
			MessageResponse: &topicspb.MessageResponse{
				Success: true,
			},
		},
	}
}

func NewCtx(msg *topicspb.ServerMessage) *Ctx {
	return &Ctx{
		id: msg.Id,
		Request: &requestImpl{
			topicName: msg.GetMessageRequest().TopicName,
			message:   msg.GetMessageRequest().Message.GetStructPayload().AsMap(),
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
