package worker

import (
	"context"
)

type Handler struct {
	Name       string
	MaxRetries int
	Handle     HandleFunc
}

type HandleFunc func(ctx context.Context, data map[string]string) error
