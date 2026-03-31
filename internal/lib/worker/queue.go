package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Queue struct {
	stream Stream
	store  Store
}

// NewQueue creates a new Queue instance with the provided Store and Stream.
// The Store is required and cannot be nil, while the Stream is optional.
// If the Stream is nil, the Queue will still function but will not publish jobs to any stream.
func NewQueue(store Store, stream Stream) *Queue {
	if store == nil {
		panic("store cannot be nil")
	}

	return &Queue{
		stream: stream,
		store:  store,
	}
}

type Task struct {
	Name string
	Data any
}

func (q *Queue) Push(ctx context.Context, task Task) error {
	data, err := json.Marshal(task.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal task data: %w", err)
	}

	job := Job{
		Name:      task.Name,
		Data:      data,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := q.store.Create(ctx, job); err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	if q.stream == nil {
		return nil
	}

	if err := q.stream.Publish(ctx, job); err != nil {
		return fmt.Errorf("failed to publish job: %w", err)
	}
	return nil
}
