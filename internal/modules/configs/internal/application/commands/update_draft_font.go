package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateDraftFontCommand struct {
	drafts models.DraftRepository
	fonts  models.FontRepository
}

func NewUpdateDraftFontCommand(drafts models.DraftRepository, fonts models.FontRepository) *UpdateDraftFontCommand {
	return &UpdateDraftFontCommand{drafts, fonts}
}

type UpdateDraftFontRequest struct {
	ConfigID primitive.ObjectID
	Family   string
	Tag      string
}

func (c *UpdateDraftFontCommand) UpdateDraftFont(ctx context.Context, req UpdateDraftFontRequest) error {
	if req.Tag == "" {
		return &models.Error{Code: "invalid_tag", Err: fmt.Errorf("tag cannot be empty")}
	}

	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
	if err != nil {
		return fmt.Errorf("failed to get draft by config ID: %w", err)
	}

	font, err := c.fonts.GetByFamily(ctx, req.Family)
	if err != nil {
		return fmt.Errorf("failed to get font by family: %w", err)
	}

	draft.Fonts[req.Tag] = font
	draft.UpdatedAt = time.Now()
	if err := c.drafts.Update(ctx, draft); err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	return nil
}
