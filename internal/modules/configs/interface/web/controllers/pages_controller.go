package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/views"
	"github.com/yaien/cultural/internal/modules/configs/library/render"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type PagesController struct {
	app *application.Application
}

func NewPagesController(app *application.Application) *PagesController {
	return &PagesController{app}
}

func (c *PagesController) Index(w http.ResponseWriter, r *http.Request) {
	_ = views.Pages().Render(r.Context(), w)
}

func (c *PagesController) List(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(config.Pages)
}

func (c *PagesController) Render(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	page, ok := config.Pages[r.PathValue("page")]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_ = render.Page(page, nil).Render(r.Context(), w)
}
