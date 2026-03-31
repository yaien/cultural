package storage

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
)

type Repository interface {
	Create(ctx context.Context, file *File) error
	Update(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id primitive.ID) (*File, error)
	GetByOrganizationIDAndName(ctx context.Context, id primitive.ID, name string) (*File, error)
	GetByOrganizationIDAndID(ctx context.Context, organizationID, id primitive.ID) (*File, error)
	GetByOrganizationID(ctx context.Context, id primitive.ID) ([]*File, error)
	RenameByOrganizationIDAndName(ctx context.Context, id primitive.ID, oldName, newName string) error
	DeleteByOrganizationIDAndName(ctx context.Context, id primitive.ID, name string) error
}
