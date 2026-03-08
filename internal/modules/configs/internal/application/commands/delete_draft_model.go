package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteDraftModelCommand struct {
	drafts models.DraftRepository
}

func NewDeleteDraftModelCommand(drafts models.DraftRepository) *DeleteDraftModelCommand {
	return &DeleteDraftModelCommand{drafts}
}

type DeleteDraftModelRequest struct {
	ConfigID primitive.ObjectID
	Type     DraftModelType
	Key      string
}

type DeleteDraftModelResponse struct {
	Draft            *models.Draft
	DefaultModelName string
	DefaultModel     any
}

func (c *DeleteDraftModelCommand) DeleteDraftModel(ctx context.Context, req DeleteDraftModelRequest) (*DeleteDraftModelResponse, error) {
	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	res := DeleteDraftModelResponse{Draft: draft}

	switch req.Type {
	case DraftPageModelType:
		_, exists := draft.Pages[req.Key]
		if !exists {
			return nil, &models.Error{Code: "not_found", Err: fmt.Errorf("page not found: %s", req.Key)}
		}

		if req.Key == models.DefaultPageName {
			return nil, &models.Error{Code: "invalid_operation", Err: fmt.Errorf("cannot delete default page")}
		}

		delete(draft.Pages, req.Key)

		res.DefaultModelName = models.DefaultPageName
		res.DefaultModel = draft.Pages[models.DefaultPageName]

	case DraftLayoutModelType:
		_, exists := draft.Layouts[req.Key]
		if !exists {
			return nil, fmt.Errorf("layout not found: %s", req.Key)
		}

		if req.Key == models.DefaultLayoutName {
			return nil, &models.Error{Code: "invalid_operation", Err: fmt.Errorf("cannot delete default layout")}
		}

		delete(draft.Layouts, req.Key)

		res.DefaultModelName = models.DefaultLayoutName
		res.DefaultModel = draft.Layouts[models.DefaultLayoutName]

	default:
		return nil, fmt.Errorf("invalid draft model type: %s", req.Type)
	}

	draft.UpdatedAt = time.Now()

	if err := c.drafts.Update(ctx, draft); err != nil {
		return nil, fmt.Errorf("failed to update draft: %w", err)
	}

	return &res, nil
}
