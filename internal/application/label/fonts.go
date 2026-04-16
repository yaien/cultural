package label

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

type Font struct {
	ID        primitive.ID      `gorm:"primaryKey;autoIncrement"`
	Family    string            `gorm:"index"`
	Subsets   []string          `gorm:"type:jsonb;serializer:json"`
	Variants  []string          `gorm:"type:jsonb;serializer:json"`
	Files     map[string]string `gorm:"type:jsonb;serializer:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Provider  string
	Category  string
	Version   string
}

type FindFontOptions struct {
	Family string
	Offset int
	Limit  int
}

type Fonts struct {
	fonts gorm.Interface[Font]
}

func NewFonts(db *gorm.DB) *Fonts {
	return &Fonts{gorm.G[Font](db)}
}

func (s *Fonts) GetByFamily(ctx context.Context, family string) (Font, error) {
	return s.fonts.Where("family = ?", family).Take(ctx)
}

func (s *Fonts) Find(ctx context.Context, opts *FindFontOptions) ([]Font, error) {
	return s.fonts.
		Where("family like ?", "%"+opts.Family+"%").
		Limit(opts.Limit).
		Offset(opts.Offset).
		Find(ctx)
}
