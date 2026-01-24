package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type CreatePageCommand struct {
	configs models.ConfigRepository
	cache   *cache.Cache[*models.Config]
}

func NewCreatePageCommand(configs models.ConfigRepository, ch *cache.Cache[*models.Config]) *CreatePageCommand {
	return &CreatePageCommand{configs, ch}
}

func (c *CreatePageCommand) CreatePage(ctx context.Context, config *models.Config, page *models.Page) error {
	if page.Name == "index" {
		return &models.Error{Code: "invalid_name", Err: errors.New("index is is a non editable page")}
	}
	if _, ok := config.Pages[page.Name]; ok {
		return &models.Error{Code: "page already exist", Err: fmt.Errorf("page %q already exists", page.Name)}
	}

	config.Pages[page.Name] = page
	config.UpdatedAt = time.Now()

	err := c.configs.Update(ctx, config)
	if err != nil {
		return fmt.Errorf("failed updating config: %w", err)
	}

	c.cache.Delete(config.Host)

	return nil
}
