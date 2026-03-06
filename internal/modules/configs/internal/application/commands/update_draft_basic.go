package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DraftModelType string

const (
	DraftPageModelType   DraftModelType = "page"
	DraftLayoutModelType DraftModelType = "layout"
	DraftEmailModelType  DraftModelType = "email"
)

type UpdateDraftBasicCommand struct {
	drafts models.DraftRepository
}

func NewUpdateDraftBasicCommand(drafts models.DraftRepository) *UpdateDraftBasicCommand {
	return &UpdateDraftBasicCommand{drafts: drafts}
}

type UpdateDraftBasicRequest struct {
	ConfigID    primitive.ObjectID
	Type        DraftModelType
	Key         string
	Name        string
	Title       string
	Description string
	Layout      string
	Subject     string
}

func (c *UpdateDraftBasicCommand) UpdateDraftBasic(ctx context.Context, req UpdateDraftBasicRequest) error {

	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	switch req.Type {
	case DraftEmailModelType:
		email, ok := draft.Emails[req.Key]
		if !ok {
			return &models.Error{Code: "not_found", Err: fmt.Errorf("email not found")}
		}

		email.Subject = req.Subject
	case DraftLayoutModelType:
		layout, ok := draft.Layouts[req.Key]
		if !ok {
			return &models.Error{Code: "not_found", Err: fmt.Errorf("layout not found")}
		}

		layout.Name = req.Name
		layout.Title = req.Title

	case DraftPageModelType:
		page, ok := draft.Pages[req.Key]
		if !ok {
			return &models.Error{Code: "not_found", Err: fmt.Errorf("page not found")}
		}

		if req.Key != models.DefaultPageName {
			page.Name = req.Name
		}

		page.Title = req.Title
		page.Description = req.Description
		page.Layout = req.Layout

	default:
		return &models.Error{Code: "invalid_type", Err: fmt.Errorf("invalid type")}
	}

	draft.UpdatedAt = time.Now()

	if err := c.drafts.Update(ctx, draft); err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	return nil

}
