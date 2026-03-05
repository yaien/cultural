package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views"
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

	var state PagesControllerStateData
	var ok bool

	state.app = c.app
	state.ctx = ctx
	state.Config = config
	state.FileURLFunc = models.FileURL
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
		views.Pages(w, r, views.Data(state), views.Template("tab_content"))
		return
	case "container":
		views.Pages(w, r, views.Data(state), views.Template("content"))
	default:
		views.Pages(w, r, views.Data(state))
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

type PagesControllerStateData struct {
	Config             *models.Config
	Draft              *models.Draft
	Selected           any
	SelectedType       string
	SelectedKey        string
	SelectedFileName   string
	SelectedFontFamily string
	FileURLFunc        models.FileURLFunc
	Section            string
	app                *application.Application
	ctx                context.Context
}

func (c PagesControllerStateData) PageIsIndex() bool {
	page, ok := c.Selected.(*models.Page)
	if !ok {
		return false
	}

	return page.Name == "index"
}

func (c PagesControllerStateData) PageUrl() string {
	page, ok := c.Selected.(*models.Page)
	if !ok {
		return ""
	}

	if page.Name == "index" {
		return c.Config.Url
	}

	return c.Config.Url + "/" + page.Name
}

func (c PagesControllerStateData) SelectedTypeOptions() any {
	switch c.Selected.(type) {
	case *models.Page:
		return c.Draft.Pages
	case *models.Layout:
		return c.Draft.Layouts
	case *models.Email:
		return c.Draft.Emails
	default:
		return nil
	}
}

func (c PagesControllerStateData) SelectedIsDeleteable() bool {
	switch sel := c.Selected.(type) {
	case *models.Page:
		return sel.Name != "index"
	case *models.Layout:
		return sel.Name != "default"
	default:
		return false
	}
}

func (c PagesControllerStateData) SelectedIsForWeb() bool {
	switch c.Selected.(type) {
	case *models.Page, *models.Layout:
		return true
	default:
		return false
	}
}

func (c PagesControllerStateData) SelectedIsAPage() bool {
	_, ok := c.Selected.(*models.Page)
	return ok
}

func (c PagesControllerStateData) SelectedIsALayout() bool {
	_, ok := c.Selected.(*models.Layout)
	return ok
}

func (c PagesControllerStateData) SelectedIsAnEmail() bool {
	_, ok := c.Selected.(*models.Email)
	return ok
}

func (c PagesControllerStateData) SelectedAttr(a any, b ...any) string {
	if len(b) == 0 {
		if selected, _ := a.(bool); selected {
			return "selected"
		}
		return ""
	}

	if a == b[0] {
		return "selected"
	}
	return ""
}

func (c PagesControllerStateData) Files() ([]*models.File, error) {
	return c.app.GetFiles(c.ctx, c.Config.OrganizationID)
}

func (c PagesControllerStateData) FileURL(name string, variant ...int) string {
	return c.FileURLFunc(name, variant...)
}

func (c PagesControllerStateData) SelectedFile() (*models.File, error) {
	file, err := c.app.GetFile(c.ctx, c.Config.OrganizationID, c.SelectedFileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}
