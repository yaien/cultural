package label

import (
	"context"
	"errors"

	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ DraftRepository = (*GormDrafts)(nil)

type GormDrafts struct {
	db *gorm.DB
}

func NewGormDrafts(db *gorm.DB) *GormDrafts {
	return &GormDrafts{db: db}
}

func (r *GormDrafts) Update(ctx context.Context, draft *Draft) error {
	return r.db.WithContext(ctx).Save(draft).Error
}

func (r *GormDrafts) GetByConfigID(ctx context.Context, configID primitive.ID) (*Draft, error) {
	var draft Draft
	err := r.db.WithContext(ctx).Where("config_id = ?", configID).First(&draft).Error
	switch {
	case err == nil:
		return &draft, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}
