package admin

import (
	"context"
	"errors"

	"github.com/yaien/cultural/internal/lib/coderror"
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
	switch {
	case err == nil:
		return &organization, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}
