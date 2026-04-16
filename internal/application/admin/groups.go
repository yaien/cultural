package admin

import (
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
