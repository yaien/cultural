package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Integration struct {
	ID             primitive.ObjectID `bson:"_id"`
	OrganizationID primitive.ObjectID `bson:"organizationId"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	Name           string             `bson:"name"`
	Data           any                `bson:"data"`
}

type GetIntegrationOptions struct {
	OrganizationID primitive.ObjectID
	Name           string
	Data           any
}

type IntegrationRepository interface {
	Create(ctx context.Context, i *Integration) error
	Update(ctx context.Context, i *Integration) error
	Get(ctx context.Context, options GetIntegrationOptions) (*Integration, error)
}
