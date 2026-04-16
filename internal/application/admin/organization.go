package admin

import (
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
)

type Organization struct {
	ID        primitive.ID `gorm:"primaryKey;autoIncrement"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
