package events

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// Topic
type Topic interface {
	Name() string
	Publish(*Event) (*Event, error)
}

type topicImpl struct {
	name string
	ec   v1.EventClient
}

func (s *topicImpl) Name() string {
	return s.name
}

func (s *topicImpl) Publish(evt *Event) (*Event, error) {
	// Convert payload to Protobuf Struct
	payloadStruct, err := structpb.NewStruct(evt.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %s", err)
	}

	r, err := s.ec.Publish(context.TODO(), &v1.EventPublishRequest{
		Topic: s.name,
		Event: &v1.NitricEvent{
			Id:          evt.ID,
			Payload:     payloadStruct,
			PayloadType: evt.PayloadType,
		},
	})

	if err != nil {
		return nil, err
	}

	// Return a reference to a new event with a populated ID
	return &Event{
		ID:          r.GetId(),
		Payload:     evt.Payload,
		PayloadType: evt.PayloadType,
	}, nil
}
