package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	Host           string             `bson:"host" json:"host"`
	Title          string             `bson:"title" json:"title"`
	Url            string             `bson:"url" json:"url"`
	Fonts          *Fonts             `bson:"fonts" json:"fonts"`
	Pages          map[string]*Page   `bson:"pages" json:"pages"`
	Colors         map[string]string  `bson:"colors" json:"colors"`
}

type ConfigRepostory interface {
	GetByHost(ctx context.Context, host string) (*Config, error)
}
