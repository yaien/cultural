package controllers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/dashboard"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/pages"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type PagesController struct {
	app *application.Application
}

func NewPagesController(app *application.Application) *PagesController {
	return &PagesController{app}
}

func (c *PagesController) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	draft, err := c.app.GetDraftByConfigID(ctx, config.ID)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed getting draft: %w", err))
		return
	}

	query := r.URL.Query()

	var ok bool

	state := &pages.State{
		Config:             config,
		Draft:              draft,
		SelectedType:       pages.SelectedType(query.Get(pages.SelectedTypeQuery)),
		SelectedKey:        query.Get(pages.SelectedKeyQuery),
		SelectedFileName:   query.Get(pages.FileQuery),
		SelectedFontFamily: query.Get(pages.FontQuery),
		Section:            query.Get(pages.SectionQuery),
		FileURL:            models.FileURL,
		Files: func() ([]*models.File, error) {
			return c.app.GetFiles(ctx, config.OrganizationID)
		},
		File: func(name string) (*models.File, error) {
			return c.app.GetFile(ctx, config.OrganizationID, name)
		},
		Fonts: func(family string, limit, offset int64) ([]*models.Font, error) {
			return c.app.GetFonts(ctx, &models.FindFontOptions{
				Family: family,
				Limit:  limit,
				Offset: int64(offset),
			})
		},
		Font: func(family string) (*models.Font, error) {
			return c.app.GetFont(ctx, family)
		},
	}

	switch state.SelectedType {
	case pages.SelectedTypeEmail:
		state.Selected, ok = draft.Emails[state.SelectedKey]
		if !ok {
			state.Selected = draft.Emails[pages.DefaultEmailName]
			state.SelectedKey = pages.DefaultEmailName
		}
	case pages.SelectedTypeLayout:
		state.Selected, ok = draft.Layouts[state.SelectedKey]
		if !ok {
			state.Selected = draft.Layouts[pages.DefaultLayoutName]
			state.SelectedKey = pages.DefaultLayoutName
		}
	default:
		state.SelectedType = pages.SelectedTypePage
		state.SelectedKey = r.URL.Query().Get(pages.SelectedKeyQuery)
		state.Selected, ok = draft.Pages[state.SelectedKey]
		if !ok {
			state.Selected = draft.Pages[pages.DefaultPageName]
			state.SelectedKey = pages.DefaultPageName
		}
	}

	switch r.Header.Get("HX-Target") {
	case pages.EditorID:
		_ = pages.Editor(state).Render(ctx, w)
		return
	case pages.ContentID:
		_ = pages.Content(state).Render(ctx, w)
	default:
		_ = pages.Page(state).Render(ctx, w)
	}

}

func (c *PagesController) Preview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)
	query := r.URL.Query()

	html, err := c.app.GetPreview(ctx, &queries.GetPreviewRequest{
		Key:    query.Get(pages.SelectedKeyQuery),
		Type:   query.Get(pages.SelectedTypeQuery),
		Config: config,
	})

	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed getting preview: %w", err))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(html))
}

func (c *PagesController) UpdateBasic(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	req := commands.UpdateDraftBasicRequest{
		ConfigID:    config.ID,
		Type:        commands.DraftModelType(r.PostForm.Get("type")),
		Key:         r.PostForm.Get("key"),
		Name:        r.PostForm.Get("name"),
		Title:       r.PostForm.Get("title"),
		Description: r.PostForm.Get("description"),
		Layout:      r.PostForm.Get("layout"),
		Subject:     r.PostForm.Get("subject"),
		OGImage:     r.PostForm.Get("og_image"),
		OGType:      r.PostForm.Get("og_type"),
	}

	if err := c.app.UpdateDraftBasic(ctx, req); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed updating basic info: %w", err))
		return
	}

	templ.Join(
		pages.Preview(req.Key, req.Type, true),
		dashboard.Toast("Cambios guardados correctamente", dashboard.Primary),
	).Render(ctx, w)

}

func (c *PagesController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	req := commands.CreateDraftModelRequest{
		ConfigID: config.ID,
		Type:     commands.DraftModelType(r.PostForm.Get("type")),
		Name:     r.PostForm.Get("name"),
		Title:    r.PostForm.Get("title"),
	}

	res, err := c.app.CreateDraftModel(ctx, req)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed creating model: %w", err))
		return
	}

	state := &pages.State{
		Config:       config,
		Draft:        res.Draft,
		SelectedType: req.Type,
		SelectedKey:  req.Name,
		Selected:     res.Model,
	}

	w.Header().Set("HX-Push-URL", fmt.Sprintf("%s?type=%s&key=%s&section=?", pages.Path, req.Type, req.Name))

	templ.Join(
		pages.Content(state),
		dashboard.Toast("Creado correctamente", dashboard.Primary),
	).Render(ctx, w)

}

func (c *PagesController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	req := commands.DeleteDraftModelRequest{
		ConfigID: config.ID,
		Type:     commands.DraftModelType(r.Form.Get("type")),
		Key:      r.Form.Get("key"),
	}

	res, err := c.app.DeleteDraftModel(ctx, req)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed deleting model: %w", err))
		return
	}

	state := &pages.State{
		Config:       config,
		Draft:        res.Draft,
		SelectedType: req.Type,
		SelectedKey:  res.DefaultModelName,
		Selected:     res.DefaultModel,
	}

	w.Header().Set("HX-Push-URL", fmt.Sprintf("%s?type=%s&key=%s&section=?", pages.Path, req.Type, res.DefaultModelName))

	templ.Join(
		pages.Content(state),
		dashboard.Toast("Eliminado correctamente", dashboard.Primary),
	).Render(ctx, w)

}
