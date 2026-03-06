package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type GetFontQuery struct {
	fonts models.FontRepository
}

func NewGetFontQuery(fonts models.FontRepository) *GetFontQuery {
	return &GetFontQuery{fonts: fonts}
}

func (q *GetFontQuery) GetFont(ctx context.Context, family string) (*models.Font, error) {
	return q.fonts.GetByFamily(ctx, family)
}
