package queues

import (
	"context"
	"fmt"

	v1 "github.com/nitrictech/go-sdk/interfaces/nitric/v1"
)

// Queue - Interface for a Queue reference
type Queue interface {
	// Name - The name of the queue
	Name() string
	// Send - Push a number of tasks to a queue
	Send([]*Task) ([]*FailedTask, error)
	// Recieve - Retrieve tasks from a queue to a maximum of the given depth
	Receive(int) ([]RecievedTask, error)
}

type queueImpl struct {
	name string
	c    v1.QueueClient
}

func (q *queueImpl) Name() string {
	return q.name
}

func (q *queueImpl) Receive(depth int) ([]RecievedTask, error) {
	if depth < 1 {
		return nil, fmt.Errorf("Depth cannot be less than 1")
	}

	r, err := q.c.Receive(context.TODO(), &v1.QueueReceiveRequest{
		Queue: q.name,
		Depth: int32(depth),
	})

	if err != nil {
		return nil, err
	}

	rts := make([]RecievedTask, len(r.GetTasks()))

	for i, task := range r.GetTasks() {
		rts[i] = &receivedTaskImpl{
			queue:   q.name,
			qc:      q.c,
			leaseId: task.GetLeaseId(),
			task:    wireToTask(task),
		}
	}

	return rts, nil
}

func (q *queueImpl) Send(tasks []*Task) ([]*FailedTask, error) {
	// Convert SDK Task objects to gRPC Task objects
	wireTasks := make([]*v1.NitricTask, len(tasks))
	for i, task := range tasks {
		wireTask, err := taskToWire(task)
		if err != nil {
			return nil, err
		}
		wireTasks[i] = wireTask
	}

	// Push the tasks to the queue
	res, err := q.c.SendBatch(context.Background(), &v1.QueueSendBatchRequest{
		Queue: q.name,
		Tasks: wireTasks,
	})
	if err != nil {
		return nil, err
	}

	// Convert the gRPC Failed Tasks to SDK Failed Task objects
	failedTasks := make([]*FailedTask, len(res.GetFailedTasks()))
	for i, failedTask := range res.GetFailedTasks() {
		failedTasks[i] = &FailedTask{
			Task:   wireToTask(failedTask.GetTask()),
			Reason: failedTask.GetMessage(),
		}
	}

	return failedTasks, nil
}
