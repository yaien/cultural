package worker

import (
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
)

type Status string

const (
	StatusPending   = Status("pending")
	StatusRunning   = Status("running")
	StatusCompleted = Status("completed")
	StatusFailed    = Status("failed")
)

type Job struct {
	ID         primitive.ID `gorm:"primaryKey;autoIncrement"`
	Name       string
	Data       []byte `gorm:"type:jsonb"`
	Status     Status
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Executions []Execution
	Retries    int
}

type Execution struct {
	StartedAt  time.Time
	FinishedAt time.Time
	FinishedIn string
	Error      string
}
