package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RoleContextKey = key("role")

type Permissions []string

type Role struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	GroupID        primitive.ObjectID `bson:"groupId,omitempty" json:"groupId,omitempty"`
	Name           string             `bson:"name" json:"name"`
	Permissions    Permissions        `bson:"permissions" json:"permissions"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type RoleRepository interface {
	GetByUserIDAndOrganizationID(ctx context.Context, userId, organizationId primitive.ObjectID) (*Role, error)
}

type Group struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	Permissions    Permissions        `bson:"permissions" json:"permissions"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type GroupRepository interface {
	GetByIDAndOrganizationID(ctx context.Context, id, organizationId primitive.ObjectID) (*Group, error)
}

func (p Permissions) Has(permission string) bool {
	for _, perm := range p {
		if perm == "*" || perm == permission {
			return true
		}
	}

	return false
}
