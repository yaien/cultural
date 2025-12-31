package commands

import (
	"context"

	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type UpdatePageCommand struct {
	configs models.ConfigRepository
	cache   *cache.Cache[*models.Config]
}

func NewUpdatePageCommand(configs models.ConfigRepository, ch *cache.Cache[*models.Config]) *UpdatePageCommand {
	return &UpdatePageCommand{configs, ch}
}

type UpdatePageRequest struct {
	Config *models.Config
	Page   *models.Page
	Path   string
}

func (c *UpdatePageCommand) UpdatePage(ctx context.Context, req *UpdatePageRequest) error {
	req.Config.Pages[req.Path] = req.Page

	err := c.configs.Update(ctx, req.Config)
	if err != nil {
		return err
	}

	c.cache.Delete(req.Config.Host)

	return nil
}
