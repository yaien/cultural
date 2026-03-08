package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateDraftColorCommand struct {
	drafts models.DraftRepository
}

func NewUpdateDraftColorCommand(drafts models.DraftRepository) *UpdateDraftColorCommand {
	return &UpdateDraftColorCommand{drafts}
}

type UpdateDraftColorRequest struct {
	ConfigID primitive.ObjectID
	ID       primitive.ObjectID
	Tag      string
	Value    string
}

func (c *UpdateDraftColorCommand) UpdateDraftColor(ctx context.Context, req *UpdateDraftColorRequest) (err error) {
	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	if req.Tag == "" {
		return &models.Error{Code: "invalid_tag", Err: fmt.Errorf("tag cannot be empty")}
	}

	if req.Value == "" {
		return &models.Error{Code: "invalid_value", Err: fmt.Errorf("value cannot be empty")}
	}

	var found bool
	for _, color := range draft.Colors {
		if color.ID.Hex() == req.ID.Hex() {
			color.Tag = req.Tag
			color.Value = req.Value
			found = true
			break
		}
	}

	if !found {
		return &models.Error{Code: "color_not_found", Err: fmt.Errorf("color with id %s not found", req.ID.Hex())}
	}

	draft.UpdatedAt = time.Now()

	if err = c.drafts.Update(ctx, draft); err != nil {
		return fmt.Errorf("failed updating draft: %w", err)
	}

	return nil

}
