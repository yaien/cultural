package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Font struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
	Provider  string             `bson:"provider"`
	Family    string             `bson:"family"`
	Category  string             `bson:"category"`
	Subsets   []string           `bson:"subsets"`
	Variants  []string           `bson:"variants"`
	Version   string             `bson:"version"`
	Files     map[string]string  `bson:"files"`
}

type FindFontOptions struct {
	Family string
	Offset int64
	Limit  int64
}

type FontRepository interface {
	Find(ctx context.Context, opts *FindFontOptions) ([]*Font, error)
	GetByFamily(ctx context.Context, family string) (*Font, error)
}
