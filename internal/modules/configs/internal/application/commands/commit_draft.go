package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type CommitDraftCommand struct {
	configs models.ConfigRepository
	drafts  models.DraftRepository
	cache   *cache.Cache[*models.Config]
}

func NewCommitDraftCommand(configs models.ConfigRepository, drafts models.DraftRepository, ch *cache.Cache[*models.Config]) *CommitDraftCommand {
	return &CommitDraftCommand{
		configs: configs,
		drafts:  drafts,
		cache:   ch,
	}
}

func (c *CommitDraftCommand) CommitDraft(ctx context.Context, config *models.Config) error {
	draft, err := c.drafts.GetByConfigID(ctx, config.ID)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	config.Colors = draft.Colors
	config.Fonts = draft.Fonts
	config.Layouts = draft.Layouts
	config.UpdatedAt = time.Now()

	if err := c.configs.Update(ctx, config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	c.cache.Delete(config.Host)
	return nil
}
