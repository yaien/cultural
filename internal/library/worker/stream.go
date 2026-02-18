package worker

import (
	"context"
	"iter"
)

type Stream interface {
	Publish(ctx context.Context, job Job) error
	Read() iter.Seq2[*Message, error]
}

type Message struct {
	Job Job
	Ack func(ctx context.Context) error
}
