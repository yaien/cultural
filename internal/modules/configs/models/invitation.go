package models

import (
	"context"
	"time"
)

type Invitation struct {
	ID             any        `bson:"_id,omitempty" json:"id"`
	OrganizationID any        `bson:"organizationId" json:"organizationId"`
	GroupID        any        `bson:"groupId,omitempty" json:"groupId,omitempty"`
	CreatorID      any        `bson:"creatorId" json:"creatorId"`
	CreatedAt      time.Time  `bson:"createdAt" json:"createdAt"`
	AcceptedAt     *time.Time `bson:"acceptedAt,omitempty" json:"acceptedAt,omitempty"`
	ExpiresAt      time.Time  `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
	Email          string     `bson:"email" json:"email"`
	Permissions    []string   `bson:"permissions" json:"permissions"`
	Name           string     `bson:"name" json:"name"`
}

type InvitationRepository interface {
	Create(ctx context.Context, invitation *Invitation) error
}
