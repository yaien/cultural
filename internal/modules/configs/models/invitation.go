package models

import (
	"context"
	"time"

	"github.com/a-h/templ"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invitation struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrganizationID  primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	CreatorID       primitive.ObjectID `bson:"creatorId,omitempty" json:"creatorId,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	AcceptedAt      *time.Time         `bson:"acceptedAt,omitempty" json:"acceptedAt,omitempty"`
	ExpiresAt       time.Time          `bson:"expiresAt" json:"expiresAt"`
	RoleGroupID     primitive.ObjectID `bson:"roleGroupId,omitempty" json:"roleGroupId,omitempty"`
	RolePermissions Permissions        `bson:"rolePermissions,omitempty" json:"rolePermissions,omitempty"`
	RoleName        string             `bson:"roleName" json:"roleName"`
	UserDisplayName string             `bson:"userDisplayName" json:"userDisplayName"`
	UserEmail       string             `bson:"userEmail" json:"userEmail"`
}

type InvitationEmailBuilder func(org *Organization, inv *Invitation, creator *User, link string) templ.Component

type InvitationRepository interface {
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Invitation, error)
	Create(ctx context.Context, invitation *Invitation) error
	Update(ctx context.Context, invitation *Invitation) error
}
