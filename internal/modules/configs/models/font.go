package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Font struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Provider  string             `bson:"provider" json:"provider"`
	Family    string             `bson:"family" json:"family"`
	Category  string             `bson:"category" json:"category"`
	Subsets   []string           `bson:"subsets" json:"subsets"`
	Variants  []string           `bson:"variants" json:"variants"`
	Version   string             `bson:"version" json:"version"`
	Files     map[string]string  `bson:"files" json:"files"`
}

type Fonts map[string]Font

type FindFontOptions struct {
	Family string
	Offset int64
	Limit  int64
}

type FontRepository interface {
	Find(ctx context.Context, opts *FindFontOptions) ([]*Font, error)
}
