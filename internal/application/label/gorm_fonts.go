package label

import (
	"context"
	"errors"
	"fmt"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ FontRepository = (*GormFonts)(nil)

type GormFonts struct {
	db *gorm.DB
}

func NewGormFonts(db *gorm.DB) *GormFonts {
	return &GormFonts{db: db}
}

func (r *GormFonts) GetByFamily(ctx context.Context, family string) (*Font, error) {
	var font Font
	err := r.db.WithContext(ctx).Where("family = ?", family).First(&font).Error
	return &font, primitive.Error(err)
}

func (r *GormFonts) Find(ctx context.Context, opts *FindFontOptions) (fonts []*Font, err error) {
	query := r.db.WithContext(ctx)

	if opts.Family != "" {
		query = query.Where("family LIKE ?", "%"+opts.Family+"%")
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}

	if opts.Offset < 0 {
		opts.Offset = 0
	}

	err = query.Offset(int(opts.Offset)).Limit(int(opts.Limit)).Find(&fonts).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed finding fonts: %w", err)
	}

	if fonts == nil {
		fonts = []*Font{}
	}

	return fonts, nil
}
