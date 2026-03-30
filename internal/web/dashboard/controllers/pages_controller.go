package controllers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/application/preview"
	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/web/dashboard/views/dashboard"
	"github.com/yaien/cultural/internal/web/dashboard/views/pages"
	"github.com/yaien/cultural/internal/web/middlewares"
)

type PagesController struct {
	drafts  *label.Drafts
	fonts   *label.Fonts
	preview *preview.Preview
	storage *storage.Storage
}

func NewPagesController(drafts *label.Drafts, fonts *label.Fonts, preview *preview.Preview, storage *storage.Storage) *PagesController {
	return &PagesController{
		drafts:  drafts,
		fonts:   fonts,
		preview: preview,
		storage: storage,
	}
}

func (c *PagesController) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	draft, err := c.drafts.GetByConfigID(ctx, config.ID)
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
		SelectedFontKey:    query.Get(pages.FontKeyQuery),
		Section:            query.Get(pages.SectionQuery),
		FileURL:            storage.FileURL,
		Files: func() ([]*storage.File, error) {
			return c.storage.GetByOrganizationID(ctx, config.OrganizationID)
		},
		File: func(name string) (*storage.File, error) {
			return c.storage.GetByOrganizationIDAndName(ctx, config.OrganizationID, name)
		},
		Fonts: func(family string, limit, offset int64) ([]*label.Font, error) {
			return c.fonts.Find(ctx, &label.FindFontOptions{
				Family: family,
				Limit:  limit,
				Offset: int64(offset),
			})
		},
		Font: func(family string) (*label.Font, error) {
			return c.fonts.GetByFamily(ctx, family)
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

	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)
	query := r.URL.Query()

	draft, err := c.drafts.GetByConfigID(ctx, config.ID)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed getting draft: %w", err))
		return
	}

	html, err := c.preview.GetHTML(ctx, &preview.GetHTMLOptions{
		Key:    query.Get(pages.SelectedKeyQuery),
		Type:   query.Get(pages.SelectedTypeQuery),
		Config: config,
		Draft:  draft,
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

	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	req := label.UpdateDraftBasicOptions{
		ConfigID:    config.ID,
		Type:        label.DraftModelType(r.PostForm.Get("type")),
		Key:         r.PostForm.Get("key"),
		Name:        r.PostForm.Get("name"),
		Title:       r.PostForm.Get("title"),
		Description: r.PostForm.Get("description"),
		Layout:      r.PostForm.Get("layout"),
		Subject:     r.PostForm.Get("subject"),
		OGImage:     r.PostForm.Get("og_image"),
		OGType:      r.PostForm.Get("og_type"),
	}

	if err := c.drafts.UpdateBasic(ctx, req); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed updating basic info: %w", err))
		return
	}

	w.Header().Set("HX-trigger", "render")
	_ = dashboard.Toast("Cambios guardados correctamente", dashboard.Primary).Render(ctx, w)

}

func (c *PagesController) UpdateSource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	req := label.UpdateDraftSourceOptions{
		ConfigID:   config.ID,
		ModelType:  label.DraftModelType(r.PostForm.Get("modelType")),
		Key:        r.PostForm.Get("key"),
		SourceType: label.DraftSourceType(r.PostForm.Get("sourceType")),
		Source:     r.PostForm.Get("source"),
	}

	if err := c.drafts.UpdateSource(ctx, &req); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed updating source: %w", err))
		return
	}

	w.Header().Set("HX-trigger", "render")
	_ = pages.Preview(req.Key, req.ModelType, true).Render(ctx, w)

}

func (c *PagesController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	req := label.CreateDraftModelOptions{
		ConfigID: config.ID,
		Type:     label.DraftModelType(r.PostForm.Get("type")),
		Name:     r.PostForm.Get("name"),
		Title:    r.PostForm.Get("title"),
	}

	res, err := c.drafts.CreateModel(ctx, req)
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

	_ = templ.Join(
		pages.Content(state),
		dashboard.Toast("Creado correctamente", dashboard.Primary),
	).Render(ctx, w)

}

func (c *PagesController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	req := label.DeleteDraftModelOptions{
		ConfigID: config.ID,
		Type:     label.DraftModelType(r.Form.Get("type")),
		Key:      r.Form.Get("key"),
	}

	res, err := c.drafts.DeleteModel(ctx, req)
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

	_ = templ.Join(
		pages.Content(state),
		dashboard.Toast("Eliminado correctamente", dashboard.Primary),
	).Render(ctx, w)

}

func (c *PagesController) CommitDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	if err := c.drafts.Commit(ctx, config); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed committing draft: %w", err))
		return
	}

	_ = dashboard.Toast("La configuración ha sido publicada correctamente", dashboard.Success).Render(ctx, w)
}
