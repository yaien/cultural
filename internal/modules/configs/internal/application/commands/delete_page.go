package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type DeletePageCommand struct {
	configs models.ConfigRepository
	cache   *cache.Cache[*models.Config]
}

func NewDeletePageCommand(configs models.ConfigRepository, ch *cache.Cache[*models.Config]) *DeletePageCommand {
	return &DeletePageCommand{
		configs: configs,
		cache:   ch,
	}
}

func (c *DeletePageCommand) DeletePage(ctx context.Context, config *models.Config, pagename string) error {
	if pagename == "index" {
		return &models.Error{Code: "invalid_page_name", Err: fmt.Errorf("index page can't be deleted")}
	}

	if _, ok := config.Pages[pagename]; !ok {
		return &models.Error{Code: "page_not_found", Err: fmt.Errorf("page %s not found", pagename)}
	}

	delete(config.Pages, pagename)
	config.UpdatedAt = time.Now()

	err := c.configs.Update(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to update config in repository: %w", err)
	}

	c.cache.Delete(config.Host)

	return nil
}
