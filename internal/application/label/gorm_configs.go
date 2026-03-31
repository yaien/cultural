package label

import (
	"context"
	"errors"

	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ ConfigRepository = (*GormConfigs)(nil)

type GormConfigs struct {
	db *gorm.DB
}

func NewGormConfigs(db *gorm.DB) *GormConfigs {
	return &GormConfigs{db: db}
}

func (r *GormConfigs) GetByHost(ctx context.Context, host string) (*Config, error) {
	var config Config
	err := r.db.WithContext(ctx).Where("host = ?", host).First(&config).Error
	switch {
	case err == nil:
		return &config, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *GormConfigs) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) (*Config, error) {
	var config Config
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).First(&config).Error
	switch {
	case err == nil:
		return &config, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *GormConfigs) Update(ctx context.Context, config *Config) error {
	return r.db.WithContext(ctx).Save(config).Error
}
