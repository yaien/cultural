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

type HandlerFunc func(ctx context.Context, data map[string]any) error

func (f HandlerFunc) Handle(ctx context.Context, data map[string]any) error {
	return f(ctx, data)
}
