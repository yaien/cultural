package storage

import (
	"github.com/yaien/cultural/internal/lib/worker"
)

var TaskName = "generate-formats"

func NewTask(file *File) worker.Task {
	return worker.Task{
		Name: TaskName,
		Data: map[string]any{
			"_id": file.ID,
		},
	}
}
