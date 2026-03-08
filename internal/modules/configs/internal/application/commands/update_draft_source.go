package commands

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateDraftSourceCommand struct {
	drafts models.DraftRepository
}

func NewUpdateDraftSourceCommand(drafts models.DraftRepository) *UpdateDraftSourceCommand {
	return &UpdateDraftSourceCommand{drafts: drafts}
}

type DraftSourceType string

const (
	DraftScriptType DraftSourceType = "script"
	DraftStylesType DraftSourceType = "styles"
	DraftBodyType   DraftSourceType = "body"
)

type UpdateDraftSourceRequest struct {
	ConfigID   primitive.ObjectID
	Source     string
	ModelType  DraftModelType
	SourceType DraftSourceType
	Key        string
}

func (c *UpdateDraftSourceCommand) UpdateDraftSource(ctx context.Context, r *UpdateDraftSourceRequest) error {
	draft, err := c.drafts.GetByConfigID(ctx, r.ConfigID)
	if err != nil {
		return models.NotFoundError(err)
	}

	switch r.ModelType {
	case DraftPageModelType:
		page, ok := draft.Pages[r.Key]
		if !ok {
			return models.NotFoundError(fmt.Errorf("page with key '%s' not found", r.Key))
		}

		switch r.SourceType {
		case DraftScriptType:
			page.Script = r.Source
		case DraftStylesType:
			page.Styles = r.Source
		case DraftBodyType:
			page.Body = r.Source
		default:
			return &models.Error{Code: "invalid_source", Err: fmt.Errorf("invalid source type: %s", r.ModelType)}
		}

	case DraftLayoutModelType:
		layout, ok := draft.Layouts[r.Key]
		if !ok {
			return models.NotFoundError(fmt.Errorf("layout with key '%s' not found", r.Key))
		}

		switch r.SourceType {
		case DraftScriptType:
			layout.Script = r.Source
		case DraftStylesType:
			layout.Styles = r.Source
		case DraftBodyType:
			layout.Body = r.Source
		default:
			return &models.Error{Code: "invalid_source", Err: fmt.Errorf("invalid source type: %s", r.ModelType)}
		}

	case DraftEmailModelType:
		email, ok := draft.Emails[r.Key]
		if !ok {
			return models.NotFoundError(fmt.Errorf("email with key '%s' not found", r.Key))
		}

		switch r.SourceType {
		case DraftBodyType:
			email.Body = r.Source
		default:
			return fmt.Errorf("invalid source type: %s", r.SourceType)
		}

	default:
		return &models.Error{Code: "invalid_model", Err: fmt.Errorf("invalid model type: %s", r.ModelType)}
	}

	if err := c.drafts.Update(ctx, draft); err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	return nil
}
