package eventclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	v1 "go.nitric.io/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type Topic interface {
	GetName() string
}

type Event struct {
	Payload     *map[string]interface{}
	PayloadType *string
	RequestId   *string
}

// Represents a Topic for event publishing. The runtime representation of a topic is provider specific.
type NitricTopic struct {
	name string
}

// GetName - returns the Nitric name of the topic
func (t *NitricTopic) GetName() string {
	return t.name
}

// String - returns the string representation of this topic
func (t *NitricTopic) String() string {
	return t.name
}

type EventClient interface {
	GetTopics() ([]Topic, error)
	Publish(opts PublishOptions) (*string, error)
}

type NitricEventClient struct {
	conn *grpc.ClientConn
	c    v1.EventingClient
}

// GetTopics - returns a slice of deployed topics in the current stack and provider.
func (e NitricEventClient) GetTopics() ([]Topic, error) {
	// Get a list of topics from the server
	res, err := e.c.GetTopics(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("an error occurred getting topics: %s", err)
	}

	// Convert the response into Topic objects
	topics := make([]Topic, 0, len(res.GetTopics()))
	for _, t := range res.GetTopics() {
		topics = append(topics, &NitricTopic{
			name: t,
		})
	}

	return topics, nil
}

type PublishOptions struct {
	TopicName *string
	Event     *Event
}

// Publish - publishes the provided event data to the specified topic.
func (e NitricEventClient) Publish(opts PublishOptions) (*string, error) {
	// Generate UUID as request id if none provided
	var requestID = opts.Event.RequestId
	if requestID == nil {
		// TODO: Pass in request id generator as an interface so it can be customized.
		uuidStr := uuid.New().String()
		if uuidStr == "" {
			return nil, fmt.Errorf("failed to generate unique request id")
		}
		requestID = &uuidStr
	}

	// Convert payload to Protobuf Struct
	payloadStruct, err := structpb.NewStruct(*opts.Event.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %s", err)
	}

	// Publish the event
	_, err = e.c.Publish(context.Background(), &v1.PublishRequest{
		TopicName: *opts.TopicName,
		Event: &v1.NitricEvent{
			RequestId:   *requestID,
			PayloadType: *opts.Event.PayloadType,
			Payload:     payloadStruct,
		},
	})

	if err != nil {
		return nil, err
	}

	return requestID, nil
}

func NewEventClient(conn *grpc.ClientConn) EventClient {
	return &NitricEventClient{
		conn: conn,
		c:    v1.NewEventingClient(conn),
	}
}

func NewWithClient(client v1.EventingClient) EventClient {
	return &NitricEventClient{
		c: client,
	}
}
