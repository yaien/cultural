package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
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

	var state pages.State
	var ok bool

	state.Config = config
	state.FileURL = models.FileURL
	state.Draft = draft
	state.SelectedType = query.Get("type")
	state.SelectedKey = query.Get("key")
	state.SelectedFileName = query.Get("file")
	state.SelectedFontFamily = query.Get("font")
	state.Section = query.Get("section")

	switch state.SelectedType {
	case "email":
		state.Selected, ok = draft.Emails[state.SelectedKey]
		if !ok {
			state.Selected = draft.Emails["invitation"]
			state.SelectedKey = "invitation"
		}
	case "layout":
		state.Selected, ok = draft.Layouts[state.SelectedKey]
		if !ok {
			state.Selected = draft.Layouts["default"]
			state.SelectedKey = "default"
		}
	default:
		state.SelectedType = "page"
		state.SelectedKey = r.URL.Query().Get("key")
		state.Selected, ok = draft.Pages[state.SelectedKey]
		if !ok {
			state.Selected = draft.Pages["index"]
			state.SelectedKey = "index"
		}
	}

	switch r.Header.Get("HX-Target") {
	case "tab-content":
		_ = pages.Section(&state).Render(ctx, w)
		return
	case "container":
		_ = pages.Content(&state).Render(ctx, w)
	default:
		_ = pages.Page(&state).Render(ctx, w)
	}

}

func (c *PagesController) Preview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)
	query := r.URL.Query()

	html, err := c.app.GetPreview(ctx, &queries.GetPreviewRequest{
		Key:    query.Get("key"),
		Type:   query.Get("type"),
		Config: config,
	})

	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed getting preview: %w", err))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(html))
}
