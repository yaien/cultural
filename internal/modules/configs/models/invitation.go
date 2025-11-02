package models

import (
	"context"
	"time"

	"github.com/a-h/templ"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invitation struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	GroupID        primitive.ObjectID `bson:"groupId,omitempty" json:"groupId,omitempty"`
	CreatorID      primitive.ObjectID `bson:"creatorId" json:"creatorId"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	AcceptedAt     *time.Time         `bson:"acceptedAt,omitempty" json:"acceptedAt,omitempty"`
	ExpiresAt      time.Time          `bson:"expiresAt" json:"expiresAt"`
	Email          string             `bson:"email" json:"email"`
	Permissions    []string           `bson:"permissions" json:"permissions"`
	Name           string             `bson:"name" json:"name"`
	DisplayName    string             `bson:"displayName,omitempty" json:"displayName,omitempty"`
}

type InvitationEmailBuilder func(org *Organization, inv *Invitation, creator *User, link string) templ.Component

type InvitationRepository interface {
	Create(ctx context.Context, invitation *Invitation) error
}
