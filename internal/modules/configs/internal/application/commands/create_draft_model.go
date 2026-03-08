package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateDraftModelCommand struct {
	drafts models.DraftRepository
}

func NewCreateDraftModelCommand(drafts models.DraftRepository) *CreateDraftModelCommand {
	return &CreateDraftModelCommand{drafts: drafts}
}

type CreateDraftModelRequest struct {
	ConfigID primitive.ObjectID
	Type     DraftModelType
	Title    string
	Name     string
}

type CreateDraftModelResponse struct {
	Draft *models.Draft
	Model any
}

func (c *CreateDraftModelCommand) CreateDraftModel(ctx context.Context, req CreateDraftModelRequest) (*CreateDraftModelResponse, error) {
	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	var res CreateDraftModelResponse
	res.Draft = draft

	switch req.Type {
	case DraftPageModelType:
		_, exists := draft.Pages[req.Name]
		if exists {
			return nil, &models.Error{Code: "already_exists", Err: fmt.Errorf("page with key '%s' already exists", req.Name)}
		}

		draft.Pages[req.Name] = &models.Page{
			Name:  req.Name,
			Title: req.Title,
		}

		res.Model = draft.Pages[req.Name]

	case DraftLayoutModelType:
		_, exists := draft.Layouts[req.Name]
		if exists {
			return nil, &models.Error{Code: "already_exists", Err: fmt.Errorf("layout with key '%s' already exists", req.Name)}
		}

		draft.Layouts[req.Name] = &models.Layout{
			Name:  req.Name,
			Title: req.Title,
		}

		res.Model = draft.Layouts[req.Name]

	default:
		return nil, &models.Error{Code: "invalid_type", Err: fmt.Errorf("invalid draft model type: %s", req.Type)}
	}

	draft.UpdatedAt = time.Now()
	if err := c.drafts.Update(ctx, draft); err != nil {
		return nil, fmt.Errorf("failed to update draft: %w", err)
	}

	return &res, nil

}
