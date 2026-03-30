package admin

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	OrganizationID primitive.ObjectID `bson:"organizationId"`
	Permissions    Permissions        `bson:"permissions"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	DeletedAt      *time.Time         `bson:"deletedAt,omitempty"`
}

type GroupRepository interface {
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Group, error)
}
