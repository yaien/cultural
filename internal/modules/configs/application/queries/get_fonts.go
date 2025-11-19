package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
)

type GetFontsQery struct {
	fonts models.FontRepository
}

func NewGetFontsQuery(fonts models.FontRepository) *GetFontsQery {
	return &GetFontsQery{fonts}
}

func (q *GetFontsQery) GetFonts(ctx context.Context, options *models.FindFontOptions) ([]*models.Font, error) {
	return q.fonts.Find(ctx, options)
}
