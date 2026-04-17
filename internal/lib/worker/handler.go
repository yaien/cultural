package worker

import (
	"context"
	"encoding/json"
	"fmt"
)

type H struct {
	Handler
	Name       string
	MaxRetries int
}

type Handler interface {
	Handle(ctx context.Context, data []byte) error
}

type HandlerFunc[T any] func(ctx context.Context, data *T) error

func (f HandlerFunc[T]) Handle(ctx context.Context, data []byte) error {
	var d T
	if err := json.Unmarshal(data, &d); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return f(ctx, &d)
}
