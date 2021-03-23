package topicclient

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

type Topic interface {
	GetName() string
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

type TopicClient interface {
	GetTopics() ([]Topic, error)
}

type NitricTopicClient struct {
	conn *grpc.ClientConn
	c    v1.TopicClient
}

// GetTopics - returns a slice of deployed topics in the current stack and provider.
func (e NitricTopicClient) GetTopics() ([]Topic, error) {
	// Get a list of topics from the server
	res, err := e.c.List(context.Background(), &v1.TopicListRequest{})
	if err != nil {
		return nil, fmt.Errorf("an error occurred getting topics: %s", err)
	}

	// Convert the response into Topic objects
	topics := make([]Topic, 0, len(res.GetTopics()))
	for _, topic := range res.GetTopics() {
		topics = append(topics, &NitricTopic{
			name: topic.GetName(),
		})
	}

	return topics, nil
}

func NewEventClient(conn *grpc.ClientConn) TopicClient {
	return &NitricTopicClient{
		conn: conn,
		c:    v1.NewTopicClient(conn),
	}
}

func NewWithClient(eventClient v1.EventClient, topicClient v1.TopicClient) TopicClient {
	return &NitricTopicClient{
		c: topicClient,
	}
}
