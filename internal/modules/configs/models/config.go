package models

import (
	"context"
	"time"
)

type Config struct {
	ID             any               `bson:"_id,omitempty" json:"id"`
	Host           string            `bson:"host" json:"host"`
	Title          string            `bson:"title" json:"title"`
	Url            string            `bson:"url" json:"url"`
	CreatedAt      time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time         `bson:"updatedAt" json:"updatedAt"`
	OrganizationID any               `bson:"organizationId" json:"organizationId"`
	Pages          map[string]*Page  `bson:"pages" json:"pages"`
	Fonts          *Fonts            `bson:"fonts" json:"fonts"`
	Colors         map[string]string `bson:"colors" json:"colors"`
}

type ConfigRepostory interface {
	GetByHost(ctx context.Context, host string) (*Config, error)
}
