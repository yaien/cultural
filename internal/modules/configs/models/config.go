package models

import (
	"context"
	"time"
)

type Config struct {
	ID             string           `bson:"_id,omitempty" json:"id"`
	Host           string           `bson:"host" json:"host"`
	Title          string           `bson:"title" json:"title"`
	Url            string           `bson:"url" json:"url"`
	Theme          string           `bson:"theme" json:"theme"`
	CreatedAt      time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time        `bson:"updatedAt" json:"updatedAt"`
	OrganizationID string           `bson:"organizationId" json:"organizationId"`
	Sites          map[string]*Site `bson:"sites" json:"sites"`
	Index          *Site            `bson:"index" json:"index"`
	Fonts          *Fonts           `bson:"fonts" json:"fonts"`
}

type ConfigRepostory interface {
	GetByHost(ctx context.Context, host string) (*Config, error)
}
