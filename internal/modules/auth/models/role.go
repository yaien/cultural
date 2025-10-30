package models

type Role struct {
	ID             any      `bson:"_id" json:"id"`
	UserID         any      `bson:"userId" json:"userId"`
	OrganizationID any      `bson:"organizationId" json:"organizationId"`
	Permissions    []string `bson:"permissions" json:"permissions"`
	Name           string   `bson:"name" json:"name"`
}
