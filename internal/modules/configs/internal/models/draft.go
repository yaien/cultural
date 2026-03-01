package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Draft struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	ConfigID  primitive.ObjectID `bson:"configId" json:"configId"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Layouts   map[string]*Layout `bson:"layouts" json:"layouts"`
	Fonts     map[string]*Font   `bson:"fonts" json:"fonts"`
	Pages     map[string]*Page   `bson:"pages" json:"pages"`
	Emails    map[string]*Email  `bson:"emails" json:"emails"`
	Colors    map[string]string  `bson:"colors" json:"colors"`
}

type DraftRepository interface {
	Update(ctx context.Context, draft *Draft) error
	GetByConfigID(ctx context.Context, configID primitive.ObjectID) (*Draft, error)
}
