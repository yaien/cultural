package controllers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/schema"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
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

type UpdateInput struct {
	Type        pages.SelectedType `schema:"type"`
	Key         string             `schema:"key"`
	Name        string             `schema:"name"`
	Title       string             `schema:"title"`
	Description string             `schema:"description"`
	Layout      string             `schema:"layout"`
	Subject     string             `schema:"subject"`
}

var decoder = schema.NewDecoder()

func (c *PagesController) UpdateBasic(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	draft, err := c.app.GetDraftByConfigID(ctx, config.ID)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed getting draft: %w", err))
		return
	}

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	var input UpdateInput
	if err := decoder.Decode(&input, r.PostForm); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed decoding form: %w", err))
		return
	}

	switch input.Type {
	case pages.SelectedTypeEmail:
		email, ok := draft.Emails[input.Key]
		if !ok {
			WriteHTMLErr(w, fmt.Errorf("email not found"))
			return
		}

		email.Subject = input.Subject
	case pages.SelectedTypeLayout:
		layout, ok := draft.Layouts[input.Key]
		if !ok {
			WriteHTMLErr(w, fmt.Errorf("layout not found"))
			return
		}

		layout.Name = input.Name
		layout.Title = input.Title

	case pages.SelectedTypePage:
		page, ok := draft.Pages[input.Key]
		if !ok {
			WriteHTMLErr(w, fmt.Errorf("page not found"))
			return
		}

		if input.Key != pages.DefaultPageName {
			page.Name = input.Name
		}

		page.Title = input.Title
		page.Description = input.Description
		page.Layout = input.Layout

	default:
		WriteHTMLErr(w, fmt.Errorf("invalid type %s", input.Type))

	}

	if err := c.app.UpdateDraft(ctx, draft); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed updating draft: %w", err))
		return
	}

	templ.Join(
		pages.Preview(input.Key, input.Type, true),
		dashboard.Toast("Cambios guardados correctamente", dashboard.Primary),
	).Render(ctx, w)

}
