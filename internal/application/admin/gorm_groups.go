package admin

import (
	"context"
	"errors"

	"github.com/yaien/cultural/internal/lib/coderror"
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
	switch {
	case err == nil:
		return &group, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}
