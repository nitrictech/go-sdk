package queues

import (
	"github.com/nitrictech/go-sdk/constants"
	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
)

// Queues - Idiomatic interface for the nitric queue service
type Queues interface {
	Queue(string) Queue
}

type queuesImpl struct {
	c v1.QueueClient
}

func (q *queuesImpl) Queue(name string) Queue {
	return &queueImpl{
		name: name,
		c:    q.c,
	}
}

// New - Construct a new Queueing Client with default options
func New() (Queues, error) {
	conn, err := grpc.Dial(
		constants.NitricAddress(),
		constants.DefaultOptions()...,
	)

	if err != nil {
		return nil, err
	}

	qClient := v1.NewQueueClient(conn)

	return &queuesImpl{
		c: qClient,
	}, nil
}
