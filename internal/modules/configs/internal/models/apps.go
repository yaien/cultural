package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppNamed interface {
	AppName() string
}

type App[T AppNamed] struct {
	ID             primitive.ObjectID `bson:"_id"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	OrganizationID primitive.ObjectID `bson:"organizationId"`
	Name           string             `bson:"name"`
	Data           T                  `bson:"data"`
}

type AppRepository[T AppNamed] interface {
	Create(context.Context, App[T]) error
	GetByOrganizationID(context.Context, primitive.ObjectID) (*App[T], error)
}
