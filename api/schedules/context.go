package schedules

import schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"

type Ctx struct {
	id       string
	Request  Request
	Response *Response
	Extras   map[string]interface{}
}

func (c *Ctx) ToClientMessage() *schedulespb.ClientMessage {
	return &schedulespb.ClientMessage{
		Id: c.id,
		Content: &schedulespb.ClientMessage_IntervalResponse{
			IntervalResponse: &schedulespb.IntervalResponse{},
		},
	}
}

func NewCtx(msg *schedulespb.ServerMessage) *Ctx {
	return &Ctx{
		id: msg.Id,
		Request: &requestImpl{
			scheduleName: msg.GetIntervalRequest().ScheduleName,
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
