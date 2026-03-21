package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Color struct {
	ID    primitive.ObjectID `bson:"_id"`
	Value string             `bson:"value"`
	Tag   string             `bson:"tag"`
}

type Colors []*Color
type Fonts map[string]*Font
type Layouts map[string]*Layout
type Pages map[string]*Page
type Emails map[string]*Email

type Draft struct {
	ID        primitive.ObjectID `bson:"_id"`
	ConfigID  primitive.ObjectID `bson:"configId"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
	Layouts   Layouts            `bson:"layouts"`
	Fonts     Fonts              `bson:"fonts"`
	Pages     Pages              `bson:"pages"`
	Emails    Emails             `bson:"emails"`
	Colors    Colors             `bson:"colors"`
}

type DraftRepository interface {
	Update(ctx context.Context, draft *Draft) error
	GetByConfigID(ctx context.Context, configID primitive.ObjectID) (*Draft, error)
}
