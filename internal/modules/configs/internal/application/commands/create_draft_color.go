package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateDraftColorCommand struct {
	drafts models.DraftRepository
}

func NewCreateDraftColorCommand(drafts models.DraftRepository) *CreateDraftColorCommand {
	return &CreateDraftColorCommand{drafts}
}

func (c *CreateDraftColorCommand) CreateDraftColor(ctx context.Context, configID primitive.ObjectID) (*models.Color, error) {
	draft, err := c.drafts.GetByConfigID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	tag, err := c.newTag(draft)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tag: %w", err)
	}

	color := &models.Color{
		ID:    primitive.NewObjectID(),
		Tag:   tag,
		Value: "#000000",
	}

	draft.Colors = append(draft.Colors, color)
	draft.UpdatedAt = time.Now()
	if err = c.drafts.Update(ctx, draft); err != nil {
		return nil, fmt.Errorf("failed updating draft: %w", err)
	}

	return color, nil
}

func (c *CreateDraftColorCommand) newTag(draft *models.Draft) (tag string, err error) {
loop:
	for index := range 100 {
		tag = fmt.Sprintf("color-%d", len(draft.Colors)+1+index)
		for _, color := range draft.Colors {
			if color.Tag == tag {
				continue loop
			}
		}
		return tag, nil
	}

	return "", fmt.Errorf("failed to generate unique key for color")
}
