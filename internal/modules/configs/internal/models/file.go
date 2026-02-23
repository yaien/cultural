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
	Preset         string             `bson:"preset" json:"preset"`
	Formats        []Format           `bson:"formats" json:"formats"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Format struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Variant     int                `bson:"variant" json:"variant"`
	Size        int64              `bson:"size" json:"size"`
	Width       int                `bson:"width" json:"width"`
	Height      int                `bson:"height" json:"height"`
	ContentType string             `bson:"contentType" json:"contentType"`
}

type FileRepository interface {
	Create(ctx context.Context, file *File) error
	Update(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*File, error)
	GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) (*File, error)
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*File, error)
	RenameByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, oldName, newName string) error
	DeleteByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) error
}

type FileURLFunc func(filename string, variant ...int) string

func NewExternalFileURLFunc(serverURL string, organizationID primitive.ObjectID) FileURLFunc {
	return func(filename string, variant ...int) string {
		if len(variant) > 0 {
			return fmt.Sprintf("%s/assets/external/%s/%s?variant=%d", serverURL, organizationID.Hex(), filename, variant[0])
		}
		return fmt.Sprintf("%s/assets/external/%s/%s", serverURL, organizationID.Hex(), filename)
	}
}
