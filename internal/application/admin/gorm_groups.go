package admin

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ GroupRepository = (*GormGroups)(nil)

type GormGroups struct {
	db *gorm.DB
}

func NewGormGroups(db *gorm.DB) *GormGroups {
	return &GormGroups{db: db}
}

func (r *GormGroups) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Group, error) {
	var group Group
	err := r.db.WithContext(ctx).Where("id = ? AND organization_id = ?", id, organizationID).First(&group).Error
	return &group, primitive.Error(err)
}
