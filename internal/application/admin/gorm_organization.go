package admin

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ OrganizationRepository = (*GormOrganizations)(nil)

type GormOrganizations struct {
	db *gorm.DB
}

func NewGormOrganizations(db *gorm.DB) *GormOrganizations {
	return &GormOrganizations{db: db}
}

func (r *GormOrganizations) GetByID(ctx context.Context, id primitive.ID) (*Organization, error) {
	var organization Organization
	err := r.db.WithContext(ctx).First(&organization, id).Error
	return &organization, primitive.Error(err)
}
