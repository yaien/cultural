package storage

import (
	"github.com/yaien/cultural/internal/lib/primitive"
	"github.com/yaien/cultural/internal/lib/worker"
)

var TaskName = "generate-formats"

type TaskData struct {
	FileID primitive.ID `json:"file_id"`
}

func NewTask(file File) worker.Task {
	return worker.Task{
		Name: TaskName,
		Data: TaskData{
			FileID: file.ID,
		},
	}
}
