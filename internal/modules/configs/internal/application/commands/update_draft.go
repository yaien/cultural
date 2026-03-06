package commands

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type UpdateDraftCommand struct {
	repo models.DraftRepository
}

func NewUpdateDraftCommand(repo models.DraftRepository) *UpdateDraftCommand {
	return &UpdateDraftCommand{repo: repo}
}

func (c *UpdateDraftCommand) UpdateDraft(ctx context.Context, draft *models.Draft) error {
	draft.UpdatedAt = time.Now()
	return c.repo.Update(ctx, draft)
}
