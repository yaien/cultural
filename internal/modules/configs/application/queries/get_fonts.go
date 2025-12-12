package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
)

type GetFontsQuery struct {
	fonts models.FontRepository
}

func NewGetFontsQuery(fonts models.FontRepository) *GetFontsQuery {
	return &GetFontsQuery{fonts}
}

func (q *GetFontsQuery) GetFonts(ctx context.Context, options *models.FindFontOptions) ([]*models.Font, error) {
	return q.fonts.Find(ctx, options)
}
