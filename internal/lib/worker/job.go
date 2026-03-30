package worker

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string

const (
	StatusPending   = Status("pending")
	StatusRunning   = Status("running")
	StatusCompleted = Status("completed")
	StatusFailed    = Status("failed")
)

type Job struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Data       map[string]any     `bson:"data"`
	Status     Status             `bson:"status"`
	CreatedAt  time.Time          `bson:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt"`
	Executions []Execution        `bson:"executions"`
	Retries    int                `bson:"retries"`
}

type Execution struct {
	StartedAt  time.Time `bson:"startedAt"`
	FinishedAt time.Time `bson:"finishedAt"`
	FinishedIn string    `bson:"finishedIn"`
	Error      string    `bson:"error,omitempty"`
}
