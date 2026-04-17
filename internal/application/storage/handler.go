package storage

import (
	"context"

	"github.com/yaien/cultural/internal/lib/worker"
)

func NewHandler(storage *Storage) worker.HandlerFunc[TaskData] {
	return func(ctx context.Context, data *TaskData) error {
		return storage.Convert(ctx, data.FileID)
	}
}
