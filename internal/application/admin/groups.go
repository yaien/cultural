package admin

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
)

type Group struct {
	ID             primitive.ID `gorm:"primaryKey;autoIncrement"`
	Name           string
	OrganizationID primitive.ID
	Permissions    Permissions `gorm:"type:jsonb;serializer:json"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type GroupRepository interface {
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Group, error)
}
