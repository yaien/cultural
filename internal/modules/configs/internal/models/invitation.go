package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invitation struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty"`
	OrganizationID  primitive.ObjectID  `bson:"organizationId"`
	CreatorID       primitive.ObjectID  `bson:"creatorId,omitempty"`
	CreatedAt       time.Time           `bson:"createdAt"`
	AcceptedAt      *time.Time          `bson:"acceptedAt,omitempty"`
	ExpiresAt       time.Time           `bson:"expiresAt"`
	RoleGroupID     *primitive.ObjectID `bson:"roleGroupId,omitempty"`
	RolePermissions Permissions         `bson:"rolePermissions,omitempty"`
	RoleName        string              `bson:"roleName"`
	UserDisplayName string              `bson:"userDisplayName"`
	UserEmail       string              `bson:"userEmail"`
}

type InvitationRepository interface {
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Invitation, error)
	Create(ctx context.Context, invitation *Invitation) error
	Update(ctx context.Context, invitation *Invitation) error
}
