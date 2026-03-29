package queries

import (
	"context"
	"fmt"
	"html/template"
	"maps"

	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type GetPreviewQuery struct {
	drafts   models.DraftRepository
	registry *models.IntegrationRegistry
}

func NewGetPreviewQuery(drafts models.DraftRepository, registry *models.IntegrationRegistry) *GetPreviewQuery {
	return &GetPreviewQuery{drafts, registry}
}

type GetPreviewRequest struct {
	Key    string
	Type   string
	Config *models.Config
}

func (q *GetPreviewQuery) GetPreview(ctx context.Context, req *GetPreviewRequest) (html string, err error) {
	draft, err := q.drafts.GetByConfigID(ctx, req.Config.ID)
	if err != nil {
		return "", fmt.Errorf("failed getting draft: %w", err)
	}

	switch req.Type {
	case "page":

		page, ok := draft.Pages[req.Key]
		if !ok {
			return "", &models.Error{Code: "page not found", Err: fmt.Errorf("page %s not found in draft pages", req.Key)}
		}

		layout, ok := draft.Layouts[page.Layout]
		if !ok {
			layout = models.DefaultLayout
		}

		return q.renderPage(ctx, page, layout, draft.Fonts, draft.Colors, req.Config)

	case "layout":
		layout, ok := draft.Layouts[req.Key]
		if !ok {
			return "", &models.Error{Code: "layout not found", Err: fmt.Errorf("layout %s not found in draft layouts", req.Key)}
		}

		page := models.EmptyPage
		return q.renderPage(ctx, page, layout, draft.Fonts, draft.Colors, req.Config)

	case "email":
		email, ok := draft.Emails[req.Key]
		if !ok {
			return "", &models.Error{Code: "email not found", Err: fmt.Errorf("email %s not found in draft emails", req.Key)}
		}
		return email.Body, nil
	default:
		return "", &models.Error{Code: "invalid type", Err: fmt.Errorf("invalid type: %s", req.Type)}
	}
}

func (q *GetPreviewQuery) renderPage(ctx context.Context, page *models.Page, layout *models.Layout, fonts models.Fonts, colors models.Colors, config *models.Config) (html string, err error) {
	funcs := template.FuncMap{}
	for _, integration := range q.registry.All() {
		if m, ok := integration.(models.IntegrationTemplateFuncMap); ok {
			fm := m.TemplateFuncMap(ctx, config)
			maps.Copy(funcs, fm)
		}
	}

	data := &models.PageData{
		Page:                page,
		Layout:              layout,
		Fonts:               fonts,
		Colors:              colors,
		FileURLFunc:         storage.FileURL,
		ExternalFileURLFunc: storage.NewExternalURLFunc(config.Url, config.OrganizationID),
		InlineStyles:        true,
		InlineScript:        true,
		Funcs:               funcs,
	}

	return models.RenderPage(data)
}
