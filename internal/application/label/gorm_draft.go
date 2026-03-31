package label

import (
	"context"

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
	return &draft, primitive.Error(err)
}
