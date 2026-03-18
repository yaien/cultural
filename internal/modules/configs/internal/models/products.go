package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organizationId"`
	Name           string             `bson:"name"`
	Slug           string             `bson:"slug"`
	Presentations  []Presentation     `bson:"presentations,omitempty"`
	Published      bool               `bson:"published"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

type Presentation struct {
	ID       primitive.ObjectID   `bson:"_id"`
	FileIDS  []primitive.ObjectID `bson:"fileIds,omitempty"`
	Name     string               `bson:"name"`
	Quantity int                  `bson:"quantity"`
	Price    float64              `bson:"price"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*Product, error)
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Product, error)
	GetBySlugAndOrganizationID(ctx context.Context, slug string, organizationID primitive.ObjectID) (*Product, error)
}
