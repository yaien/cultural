package worker

import (
	"context"
	"iter"
)

var _ Stream = (*MemoryStream)(nil)

type MemoryStream struct {
	channel chan Job
}

func NewMemoryStream() *MemoryStream {
	return &MemoryStream{
		channel: make(chan Job, 100),
	}
}

func (m *MemoryStream) Publish(ctx context.Context, job Job) error {
	go func() {
		select {
		case <-ctx.Done():
		case m.channel <- job:
		}
	}()

	return nil
}

func (m *MemoryStream) Read(ctx context.Context) iter.Seq2[*Message, error] {
	return func(yield func(*Message, error) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case job := <-m.channel:
				m := &Message{
					Job: job,
					Ack: func(ctx context.Context) error { return nil },
				}

				if !yield(m, nil) {
					break
				}
			}
		}
	}
}
