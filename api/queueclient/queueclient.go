package queueclient

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

type Task struct {
	// ID - Unique ID for the task
	ID string
	// LeaseID - (Read-Only) LeaseID that can be used to complete this task
	LeaseID string
	// PayloadType - Deserialization hint for interprocess communication
	PayloadType string
	// Payload - The payload to include in this task
	Payload map[string]interface{}
}

// FailedTask - A wrapper for returning errors when tasks fail to enqueue
type FailedTask struct {
	Task    *Task
	Message string
}

// SendBatchResponse - Response for SendBatch API call
type SendBatchResult struct {
	FailedTasks []*FailedTask
}

type SendOptions struct {
	Queue string
	Task  *Task
}

type SendBatchOptions struct {
	Queue string
	Tasks []*Task
}

type RecieveOptions struct {
	Queue string
	Depth int
}

type CompleteOptions struct {
	Queue string
	Task  *Task
}

type QueueRecieveResult struct {
}

type SendResult struct{}
type RecieveResult struct {
	Tasks []*Task
}
type CompleteResult struct{}

type QueueClient interface {
	Send(opts *SendOptions) (*SendResult, error)
	SendBatch(opts *SendBatchOptions) (*SendBatchResult, error)
	Receive(opts *RecieveOptions) (*RecieveResult, error)
	Complete(opts *CompleteOptions) (*CompleteResult, error)
}

type NitricQueueClient struct {
	conn *grpc.ClientConn
	c    v1.QueueClient
}

func taskToWire(task *Task) (*v1.NitricTask, error) {
	// Convert payload to Protobuf Struct
	payloadStruct, err := structpb.NewStruct(task.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %s", err)
	}

	return &v1.NitricTask{
		Id:          task.ID,
		PayloadType: task.PayloadType,
		Payload:     payloadStruct,
	}, nil
}

func wireToTask(task *v1.NitricTask) *Task {
	return &Task{
		ID:          task.GetId(),
		PayloadType: task.GetPayloadType(),
		Payload:     task.GetPayload().AsMap(),
		LeaseID:     task.GetLeaseId(),
	}
}

// Send - Sends a single task to a queue to be processed asynchronously by other services
// queueName is the nitric name of the queue (defined in the nitric.yaml file)
func (q NitricQueueClient) Send(opts *SendOptions) (*SendResult, error) {
	var finalErr error
	if wireTask, err := taskToWire(opts.Task); err == nil {
		if _, err = q.c.Send(context.TODO(), &v1.QueueSendRequest{
			Queue: opts.Queue,
			Task:  wireTask,
		}); err == nil {
			return &SendResult{}, nil
		}

		finalErr = err
	}

	return nil, finalErr
}

// SendBatch - publishes multiple tasks to a queue to be processed asynchronously by other services
// queueName should be the Nitric name of the queue. This will be automatically resolved to the provider specific
// queue identifier.
func (q NitricQueueClient) SendBatch(opts *SendBatchOptions) (*SendBatchResult, error) {
	// Convert SDK Task objects to gRPC Task objects
	wireTasks := make([]*v1.NitricTask, len(opts.Tasks))
	for i, task := range opts.Tasks {
		wireTask, err := taskToWire(task)
		if err != nil {
			return nil, err
		}
		wireTasks[i] = wireTask
	}

	// Push the tasks to the queue
	res, err := q.c.SendBatch(context.Background(), &v1.QueueSendBatchRequest{
		Queue: opts.Queue,
		Tasks: wireTasks,
	})
	if err != nil {
		return nil, err
	}

	// Convert the gRPC Failed Tasks to SDK Failed Task objects
	failedTasks := make([]*FailedTask, len(res.GetFailedTasks()))
	for i, failedTask := range res.GetFailedTasks() {
		failedTasks[i] = &FailedTask{
			Task:    wireToTask(failedTask.GetTask()),
			Message: failedTask.GetMessage(),
		}
	}

	return &SendBatchResult{FailedTasks: failedTasks}, nil
}

// Receive - retrieve tasks from the specifed queue. The items returned are contained in a QueueItem
// which provides context for the source queue and the lease on the tasks.
// queue items must be completed using Complete or they will be distributed again or forwarded to a dead letter queue.
func (q NitricQueueClient) Receive(opts *RecieveOptions) (*RecieveResult, error) {
	// Set minimum depth to 1.
	var depth = 1
	if opts.Depth > 0 {
		depth = opts.Depth
	}

	// receieve up to the requested depth off of the queue
	res, err := q.c.Receive(context.Background(), &v1.QueueReceiveRequest{
		Queue: opts.Queue,
		Depth: int32(depth),
	})
	if err != nil {
		return nil, err
	}

	// Convert the items to SDK QueueItem objects
	tasks := make([]*Task, len(res.GetTasks()))
	for i, item := range res.GetTasks() {
		tasks[i] = wireToTask(item)
	}

	return &RecieveResult{
		Tasks: tasks,
	}, nil
}

// Complete - marks a queue item as successfully completed and removes it from the queue.
//
// All items retrieved through Pop must be Completed or Released so they're not reprocessed or sent to a dead letter queue.
func (q NitricQueueClient) Complete(opts *CompleteOptions) (*CompleteResult, error) {
	_, err := q.c.Complete(context.Background(), &v1.QueueCompleteRequest{
		Queue:   opts.Queue,
		LeaseId: opts.Task.LeaseID,
	})

	return &CompleteResult{}, err
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
