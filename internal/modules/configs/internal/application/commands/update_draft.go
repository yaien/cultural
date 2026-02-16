package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateDraftCommand struct {
	repo models.DraftRepository
}

func NewUpdateDraftCommand(repo models.DraftRepository) *UpdateDraftCommand {
	return &UpdateDraftCommand{repo: repo}
}

type UpdateDraftRequest struct {
	ConfigID primitive.ObjectID
	Fonts    map[string]*models.Font
	Pages    map[string]*models.Page
	Layouts  map[string]*models.Page
	Emails   map[string]*models.Email
	Colors   map[string]string
}

func (c *UpdateDraftCommand) UpdateDraft(ctx context.Context, r *UpdateDraftRequest) error {
	draft, err := c.repo.GetByConfigID(ctx, r.ConfigID)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	draft.Fonts = r.Fonts
	draft.Pages = r.Pages
	draft.Emails = r.Emails
	draft.Colors = r.Colors
	draft.Layouts = r.Layouts
	draft.UpdatedAt = time.Now()

	return c.repo.Update(ctx, draft)
}
