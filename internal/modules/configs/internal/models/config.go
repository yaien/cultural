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
	Email          string             `bson:"email" json:"email"`
	Fonts          Fonts              `bson:"fonts" json:"fonts"`
	Pages          Pages              `bson:"pages" json:"pages"`
	Layouts        Layouts            `bson:"layouts" json:"layouts"`
	Emails         Emails             `bson:"emails" json:"emails"`
	Colors         Colors             `bson:"colors" json:"colors"`
}

type ConfigRepository interface {
	GetByHost(ctx context.Context, host string) (*Config, error)
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (*Config, error)
	Update(ctx context.Context, config *Config) error
}
