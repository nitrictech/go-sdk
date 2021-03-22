package queueclient

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"github.com/nitrictech/go-sdk/v1/eventclient"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type FailedEvent struct {
	event   eventclient.Event
	message string
}

type SendBatchResponse struct {
	failedEvents []FailedEvent
}

type QueueItem struct {
	event   eventclient.Event
	leaseId string
	queue   string
}

type QueueClient interface {
	Send(queueName string, event eventclient.Event) error
	SendBatch(queueName string, events []eventclient.Event) (*SendBatchResponse, error)
	Receive(queueName string, depth int) ([]QueueItem, error)
	Complete(item QueueItem) error
}

type NitricQueueClient struct {
	conn *grpc.ClientConn
	c    v1.QueueClient
}

func eventToWire(event eventclient.Event) (*v1.NitricEvent, error) {
	// Convert payload to Protobuf Struct
	payloadStruct, err := structpb.NewStruct(*event.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %s", err)
	}

	return &v1.NitricEvent{
		RequestId:   *event.RequestId,
		PayloadType: *event.PayloadType,
		Payload:     payloadStruct,
	}, nil
}

func wireToEvent(event *v1.NitricEvent) eventclient.Event {
	payload := event.Payload.AsMap()
	return eventclient.Event{
		RequestId:   &event.RequestId,
		PayloadType: &event.PayloadType,
		Payload:     &payload,
	}
}

// Send - Sends a single event to a queue to be processed asynchronously by other services
// queueName is the nitric name of the queue (defined in the nitric.yaml file)
func (q NitricQueueClient) Send(queueName string, event eventclient.Event) error {
	wireEvent, err := eventToWire(event)

	if err == nil {
		_, err = q.c.Send(context.TODO(), &v1.QueueSendRequest{
			Queue: queueName,
			Event: wireEvent,
		})
	}

	return err
}

// SendBatch - publishes multiple events to a queue to be processed asynchronously by other services
// queueName should be the Nitric name of the queue. This will be automatically resolved to the provider specific
// queue identifier.
func (q NitricQueueClient) SendBatch(queueName string, events []eventclient.Event) (*SendBatchResponse, error) {
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
	res, err := q.c.SendBatch(context.Background(), &v1.QueueSendBatchRequest{
		Queue:  queueName,
		Events: wireEvents,
	})
	if err != nil {
		return nil, err
	}

	// Convert the gRPC Failed Events to SDK Failed Event objects
	failedEvents := make([]FailedEvent, len(res.GetFailedEvents()))
	for i, failedEvent := range res.GetFailedEvents() {
		failedEvents[i] = FailedEvent{
			event:   wireToEvent(failedEvent.GetEvent()),
			message: failedEvent.GetMessage(),
		}
	}

	return &SendBatchResponse{failedEvents: failedEvents}, nil
}

// Receive - retrieve events from the specifed queue. The items returned are contained in a QueueItem
// which provides context for the source queue and the lease on the event.
// queue items must be completed using Complete or they will be distributed again or forwarded to a dead letter queue.
func (q NitricQueueClient) Receive(queueName string, depth int) ([]QueueItem, error) {
	// Set minimum depth to 1.
	if depth < 1 {
		depth = 1
	}

	// receieve up to the requested depth off of the queue
	res, err := q.c.Receive(context.Background(), &v1.QueueReceiveRequest{
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
			queue:   queueName,
		}
	}

	return queueItems, nil
}

// Complete - marks a queue item as successfully completed and removes it from the queue.
//
// All items retrieved through Pop must be Completed or Released so they're not reprocessed or sent to a dead letter queue.
func (q NitricQueueClient) Complete(item QueueItem) error {
	_, err := q.c.Complete(context.Background(), &v1.QueueCompleteRequest{
		Queue:   item.queue,
		LeaseId: item.leaseId,
	})

	return err
}

func NewQueueClient(conn *grpc.ClientConn) QueueClient {
	return &NitricQueueClient{
		conn: conn,
		c:    v1.NewQueueClient(conn),
	}
}

func NewWithClient(client v1.QueueClient) QueueClient {
	return &NitricQueueClient{
		c: client,
	}
}
