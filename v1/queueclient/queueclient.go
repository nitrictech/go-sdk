package queueclient

import (
	"context"
	"fmt"
	v1 "go.nitric.io/go-sdk/interfaces/nitric/v1"
	"go.nitric.io/go-sdk/v1/eventclient"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type FailedEvent struct {
	event eventclient.Event
	message string
}

type PushResponse struct {
	failedEvents []FailedEvent
}

type QueueItem struct {
	event eventclient.Event
	leaseId string
	queue string
}

type QueueClient interface {
	Push (queueName string, events []eventclient.Event) (*PushResponse, error)
	Pop (queueName string, depth int) ([]QueueItem, error)
	Complete (item QueueItem) error
}

type NitricQueueClient struct {
	conn *grpc.ClientConn
	c v1.QueueClient
}

func eventToWire(event eventclient.Event) (*v1.NitricEvent, error)  {
	// Convert payload to Protobuf Struct
	payloadStruct, err := structpb.NewStruct(*event.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %s", err)
	}

	return &v1.NitricEvent{
		RequestId: *event.RequestId,
		PayloadType: *event.PayloadType,
		Payload: payloadStruct,
	}, nil
}

func wireToEvent(event *v1.NitricEvent) eventclient.Event  {
	payload := event.Payload.AsMap()
	return eventclient.Event{
		RequestId:   &event.RequestId,
		PayloadType: &event.PayloadType,
		Payload:     &payload,
	}
}

// Push - publishes events to a queue to be processed asynchronously by other services
// queueName should be the Nitric name of the queue. This will be automatically resolved to the provider specific
// queue identifier.
func (q NitricQueueClient) Push(queueName string, events []eventclient.Event) (*PushResponse, error)  {
	// Convert SDK Event objects to gRPC Event objects
	wireEvents := make([]*v1.NitricEvent, len(events))
	for i, event := range events {
		wireEvent, err := eventToWire(event)
		if err != nil {
			return nil, err
		}
		wireEvents[i] = wireEvent
	}

	// Push the events to the queue
	res, err := q.c.Push(context.Background(), &v1.PushRequest{
		Queue:  queueName,
		Events: wireEvents,
	})
	if err != nil {
		return nil, err
	}

	// Convert the gRPC Failed Events to SDK Failed Event objects
	failedEvents := make([]FailedEvent, len(res.GetFailedMessages()))
	for i, failedEvent := range res.GetFailedMessages() {
		failedEvents[i] = FailedEvent{
			event: wireToEvent(failedEvent.GetEvent()),
			message: failedEvent.GetMessage(),
		}
	}

	return &PushResponse{failedEvents: failedEvents}, nil
}

// Pop - retrieve events from the specifed queue. The items returned are contained in a QueueItem
// which provides context for the source queue and the lease on the event.
// queue items must be completed using Complete or they will be distributed again or forwarded to a dead letter queue.
func (q NitricQueueClient) Pop(queueName string, depth int) ([]QueueItem, error)  {
	// Set minimum depth to 1.
	if depth < 1 {
		depth = 1
	}

	// Pop the requested off the queue
	res, err := q.c.Pop(context.Background(), &v1.PopRequest{
		Queue: queueName,
		Depth: int32(depth),
	})
	if err != nil {
		return nil, err
	}

	// Convert the items to SDK QueueItem objects
	queueItems := make([]QueueItem, len(res.GetItems()))
	for i, item := range res.GetItems() {
		queueItems[i] = QueueItem{
			event:   wireToEvent(item.GetEvent()),
			leaseId: item.GetLeaseId(),
			queue: queueName,
		}
	}

	return queueItems, nil
}

// Complete - marks a queue item as successfully completed and removes it from the queue.
//
// All items retrieved through Pop must be Completed or Released so they're not reprocessed or sent to a dead letter queue.
func (q NitricQueueClient) Complete(item QueueItem) error {
	_, err := q.c.Complete(context.Background(), &v1.CompleteRequest{
		Queue: item.queue,
		LeaseId: item.leaseId,
	})

	return err
}

// Close - closes the connection to the membrane server
// no need to call close if the connect is to remain open for the lifetime of the application.
func (q NitricQueueClient) Close() error {
	return q.conn.Close()
}

func NewQueueClient() (QueueClient, error) {
	// Connect to the gRPC Membrane Server
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to Membrane gRPC server: %s", err)
	}

	return &NitricQueueClient{
		conn: conn,
		c: v1.NewQueueClient(conn),
	}, nil
}