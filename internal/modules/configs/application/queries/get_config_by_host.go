package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"github.com/yaien/cultural/internal/shared"
)

type GetConfigByHostQuery struct {
	repo  models.ConfigRepostory
	cache *shared.Cache[*models.Config]
}

func NewGetConfigByHostQuery(repo models.ConfigRepostory, cache *shared.Cache[*models.Config]) *GetConfigByHostQuery {
	return &GetConfigByHostQuery{
		repo:  repo,
		cache: cache,
	}
}

func (q *GetConfigByHostQuery) GetConfigByHost(ctx context.Context, host string) (*models.Config, error) {
	if config, found := q.cache.Get(host); found {
		return config, nil
	}

	config, err := q.repo.GetByHost(ctx, host)
	if err != nil {
		return nil, err
	}

	q.cache.Set(host, config)

	return config, nil
}
