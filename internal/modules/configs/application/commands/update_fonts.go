package commands

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type UpdateFontsCommand struct {
	configs models.ConfigRepository
	cache   *cache.Cache[*models.Config]
}

func NewUpdateFontsCommand(configs models.ConfigRepository, ch *cache.Cache[*models.Config]) *UpdateFontsCommand {
	return &UpdateFontsCommand{configs, ch}
}

func (c *UpdateFontsCommand) UpdateFonts(ctx context.Context, config models.Config, fonts models.Fonts) error {
	config.Fonts = fonts

	err := c.configs.Update(ctx, &config)
	if err != nil {
		return err
	}

	c.cache.Delete(config.Host)

	return nil
}
