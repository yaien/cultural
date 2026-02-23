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
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Data       map[string]any     `bson:"data" json:"data"`
	Status     Status             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
	Executions []Execution        `bson:"executions" json:"executions"`
	Retries    int                `bson:"retries" json:"retries"`
}

type Execution struct {
	StartedAt  time.Time `bson:"startedAt" json:"startedAt"`
	FinishedAt time.Time `bson:"finishedAt" json:"finishedAt"`
	FinishedIn string    `bson:"finishedIn" json:"finishedIn"`
	Error      string    `bson:"error,omitempty" json:"error,omitempty"`
}
