package models

import "time"

type Role struct {
	ID             any        `bson:"_id,omitempty" json:"id"`
	UserID         any        `bson:"userId" json:"userId"`
	OrganizationID any        `bson:"organizationId" json:"organizationId"`
	GroupID        any        `bson:"groupId,omitempty" json:"groupId,omitempty"`
	Name           string     `bson:"name" json:"name"`
	Permissions    []string   `bson:"permissions" json:"permissions"`
	CreatedAt      time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

type Group struct {
	ID             any        `bson:"_id,omitempty" json:"id"`
	Name           string     `bson:"name" json:"name"`
	OrganizationID any        `bson:"organizationId" json:"organizationId"`
	Permissions    []string   `bson:"permissions" json:"permissions"`
	CreatedAt      time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time  `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}
