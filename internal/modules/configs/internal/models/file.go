package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	Name           string             `bson:"name" json:"name"`
	ContentType    string             `bson:"contentType" json:"contentType"`
	Formats        map[int]Format     `bson:"formats" json:"formats"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Format struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Variant int                `bson:"variant" json:"variant"`
	Size    int64              `bson:"size" json:"size"`
	Width   int                `bson:"width" json:"width"`
	Height  int                `bson:"height" json:"height"`
}

type FileRepository interface {
	Create(ctx context.Context, file *File) error
	Get(ctx context.Context, organizationID primitive.ObjectID, name string) (*File, error)
	Rename(ctx context.Context, organizationID primitive.ObjectID, oldName, newName string) error
	Delete(ctx context.Context, organizationID primitive.ObjectID, name string) error
	List(ctx context.Context, organizationID primitive.ObjectID) ([]*File, error)
}

type FileURLFunc func(filename string) string

func NewExternalFileURLFunc(serverURL string, organizationID primitive.ObjectID) FileURLFunc {
	return func(filename string) string {
		return fmt.Sprintf("%s/assets/external/%s/%s", serverURL, organizationID.Hex(), filename)
	}
}
