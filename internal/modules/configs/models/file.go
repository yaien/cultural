package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	Name           string             `bson:"name" json:"name"`
	Size           int64              `bson:"size" json:"size"`
	MimeType       string             `bson:"mimeType" json:"mimeType"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type FileRepository interface {
	Create(ctx context.Context, file *File) error
	Get(ctx context.Context, organizationID primitive.ObjectID, name string) (*File, error)
	Rename(ctx context.Context, organizationID primitive.ObjectID, oldName, newName string) error
	Delete(ctx context.Context, organizationID primitive.ObjectID, name string) error
	List(ctx context.Context, organizationID primitive.ObjectID) ([]*File, error)
}
