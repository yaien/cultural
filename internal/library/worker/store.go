package worker

import "context"

type Store interface {
	Fetch(ctx context.Context) ([]Job, error)
	Update(ctx context.Context, job Job) error
	Create(ctx context.Context, job Job) error
}
