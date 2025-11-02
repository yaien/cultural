package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type GetConfigByHostQuery struct {
	repo  models.ConfigRepository
	cache *cache.Cache[*models.Config]
}

func NewGetConfigByHostQuery(repo models.ConfigRepository, ch *cache.Cache[*models.Config]) *GetConfigByHostQuery {
	return &GetConfigByHostQuery{
		repo:  repo,
		cache: ch,
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
