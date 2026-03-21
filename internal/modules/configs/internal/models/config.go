package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organizationId"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	Host           string             `bson:"host"`
	Title          string             `bson:"title"`
	Url            string             `bson:"url"`
	Email          string             `bson:"email"`
	Fonts          Fonts              `bson:"fonts"`
	Pages          Pages              `bson:"pages"`
	Layouts        Layouts            `bson:"layouts"`
	Emails         Emails             `bson:"emails"`
	Colors         Colors             `bson:"colors"`
}

type ConfigRepository interface {
	GetByHost(ctx context.Context, host string) (*Config, error)
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (*Config, error)
	Update(ctx context.Context, config *Config) error
}
