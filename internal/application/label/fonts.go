package label

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
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
	Offset int64
	Limit  int64
}

type FontRepository interface {
	Find(ctx context.Context, opts *FindFontOptions) ([]*Font, error)
	GetByFamily(ctx context.Context, family string) (*Font, error)
}

type Fonts struct {
	fonts FontRepository
}

func NewFonts(fonts FontRepository) *Fonts {
	return &Fonts{fonts: fonts}
}

func (s *Fonts) GetByFamily(ctx context.Context, family string) (*Font, error) {
	return s.fonts.GetByFamily(ctx, family)
}

func (s *Fonts) Find(ctx context.Context, opts *FindFontOptions) ([]*Font, error) {
	return s.fonts.Find(ctx, opts)
}
