package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteDraftColorCommand struct {
	drafts models.DraftRepository
}

func NewDeleteDraftColorCommand(drafts models.DraftRepository) *DeleteDraftColorCommand {
	return &DeleteDraftColorCommand{drafts: drafts}
}

func (c *DeleteDraftColorCommand) DeleteDraftColor(ctx context.Context, configID, id primitive.ObjectID) error {
	draft, err := c.drafts.GetByConfigID(ctx, configID)
	if err != nil {
		return err
	}

	var deleted bool
	for i, color := range draft.Colors {
		if color.ID == id {
			draft.Colors = append(draft.Colors[:i], draft.Colors[i+1:]...)
			deleted = true
			break
		}
	}

	if !deleted {
		return &models.Error{Code: "not_found", Err: fmt.Errorf("no color found with id %s", id.Hex())}
	}

	draft.UpdatedAt = time.Now()
	return c.drafts.Update(ctx, draft)
}
