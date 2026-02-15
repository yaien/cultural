package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
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
	views.Pages(w, r)
}

func (c *PagesController) List(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	WriteJSON(w, config.Pages)
}

func (c *PagesController) Create(w http.ResponseWriter, r *http.Request) {

	var page models.Page

	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		WriteJSONErr(w, models.DecodeError(err))
		return
	}

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	if err := c.app.CreatePage(r.Context(), config, &page); err != nil {
		WriteJSONErr(w, fmt.Errorf("failed creating page: %w", err))
		return
	}

	WriteJSON(w, map[string]any{"path": page.Name})
}

func (c *PagesController) Update(w http.ResponseWriter, r *http.Request) {
	var page models.Page
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		WriteJSONErr(w, models.DecodeError(err))
		return
	}

	ctx := r.Context()
	path := r.PathValue("page")
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	request := &commands.UpdatePageRequest{
		Page:   &page,
		Config: config,
		Path:   path,
	}

	if err := c.app.UpdatePage(ctx, request); err != nil {
		WriteJSONErr(w, fmt.Errorf("failed updating page: %w", err))
		return
	}

	WriteJSON(w, map[string]any{"path": path})

}

func (c *PagesController) Delete(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("page")
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	if err := c.app.DeletePage(r.Context(), config, path); err != nil {
		WriteJSONErr(w, fmt.Errorf("failed deleting page: %w", err))
		return
	}

	WriteJSON(w, map[string]any{"path": path})
}
