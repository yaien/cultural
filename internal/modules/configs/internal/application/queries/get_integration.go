package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type GetIntegrationQuery struct {
	integrations models.IntegrationRepository
}

func NewGetIntegrationQuery(repo models.IntegrationRepository) *GetIntegrationQuery {
	return &GetIntegrationQuery{
		integrations: repo,
	}
}

type GetIntegrationOptions = models.GetIntegrationOptions

func (q *GetIntegrationQuery) GetIntegration(ctx context.Context, options models.GetIntegrationOptions) (*models.Integration, error) {
	return q.integrations.Get(ctx, options)
}
