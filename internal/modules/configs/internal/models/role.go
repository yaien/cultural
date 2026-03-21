package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permissions []string

type Role struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	UserID         primitive.ObjectID  `bson:"userId"`
	UserEmail      string              `bson:"userEmail"`
	UserName       string              `bson:"userName"`
	UserAvatarUrl  string              `bson:"userAvatarUrl"`
	OrganizationID primitive.ObjectID  `bson:"organizationId"`
	GroupID        *primitive.ObjectID `bson:"groupId,omitempty"`
	Name           string              `bson:"name"`
	Permissions    Permissions         `bson:"permissions"`
	CreatedAt      time.Time           `bson:"createdAt"`
	UpdatedAt      time.Time           `bson:"updatedAt"`
	DeletedAt      *time.Time          `bson:"deletedAt,omitempty"`
}

type RoleRepository interface {
	CountAdminsByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (int64, error)
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Role, error)
	GetByUserIDAndOrganizationID(ctx context.Context, userId, organizationID primitive.ObjectID) (*Role, error)
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, role *Role) error
}

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

func (p Permissions) Has(permission string) bool {
	for _, perm := range p {
		if perm == "*" || perm == permission {
			return true
		}
	}

	return false
}
