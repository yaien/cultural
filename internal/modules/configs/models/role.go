package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	GroupID        primitive.ObjectID `bson:"groupId,omitempty" json:"groupId,omitempty"`
	Name           string             `bson:"name" json:"name"`
	Permissions    []string           `bson:"permissions" json:"permissions"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type Group struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	Permissions    []string           `bson:"permissions" json:"permissions"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}
