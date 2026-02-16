package commands

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type UpdateColorsCommand struct {
	configs models.ConfigRepository
	cache   *cache.Cache[*models.Config]
}

func NewUpdateColorsCommand(configs models.ConfigRepository, ch *cache.Cache[*models.Config]) *UpdateColorsCommand {
	return &UpdateColorsCommand{configs, ch}
}

func (c *UpdateColorsCommand) UpdateColors(ctx context.Context, config *models.Config, colors map[string]string) error {
	config.Colors = colors
	config.UpdatedAt = time.Now()

	err := c.configs.Update(ctx, config)
	if err != nil {
		return err
	}

	c.cache.Delete(config.Host)

	return nil
}
