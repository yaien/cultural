package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(ctx context.Context, file *File) error
	Update(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*File, error)
	GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) (*File, error)
	GetByOrganizationIDAndID(ctx context.Context, organizationID, id primitive.ObjectID) (*File, error)
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*File, error)
	RenameByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, oldName, newName string) error
	DeleteByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) error
}
