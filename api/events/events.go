package events

import (
	"context"

	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

// Events
type Events interface {
	// Topic - Retrieve a Topic reference
	Topic(name string) Topic
	// Topics - Retrievs a list of available Topic references
	Topics() ([]Topic, error)
}

type eventsImpl struct {
	ec v1.EventClient
	tc v1.TopicClient
}

func (s *eventsImpl) Topic(name string) Topic {
	// Just return the straight topic reference
	// we can fail if the topic does not exist
	return &topicImpl{
		name: name,
		ec:   s.ec,
	}
}

func (s *eventsImpl) Topics() ([]Topic, error) {
	r, err := s.tc.List(context.TODO(), &v1.TopicListRequest{})

	if err != nil {
		return nil, err
	}

	ts := make([]Topic, 0)
	for _, topic := range r.GetTopics() {
		ts = append(ts, s.Topic(topic.GetName()))
	}

	return ts, nil
}

// New - Construct a new Eventing Client with default options
func New() (Events, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)

	if err != nil {
		return nil, err
	}

	ec := v1.NewEventClient(conn)
	tc := v1.NewTopicClient(conn)

	return &eventsImpl{
		ec: ec,
		tc: tc,
	}, nil
}
