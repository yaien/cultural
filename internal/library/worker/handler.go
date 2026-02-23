package worker

import (
	"context"
)

type H struct {
	Handler
	Name       string
	MaxRetries int
}

type Handler interface {
	Handle(ctx context.Context, data map[string]any) error
}
