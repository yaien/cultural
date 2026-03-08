package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Color struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Value string             `bson:"value" json:"value"`
	Tag   string             `bson:"tag" json:"tag"`
}

type Colors []*Color
type Fonts map[string]*Font
type Layouts map[string]*Layout
type Pages map[string]*Page
type Emails map[string]*Email

type Draft struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	ConfigID  primitive.ObjectID `bson:"configId" json:"configId"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Layouts   Layouts            `bson:"layouts" json:"layouts"`
	Fonts     Fonts              `bson:"fonts" json:"fonts"`
	Pages     Pages              `bson:"pages" json:"pages"`
	Emails    Emails             `bson:"emails" json:"emails"`
	Colors    Colors             `bson:"colors" json:"colors"`
}

type DraftRepository interface {
	Update(ctx context.Context, draft *Draft) error
	GetByConfigID(ctx context.Context, configID primitive.ObjectID) (*Draft, error)
}
