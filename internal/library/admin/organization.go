package admin

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type OrganizationRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*Organization, error)
}
