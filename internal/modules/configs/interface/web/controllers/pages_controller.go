package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/views"
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

func (c *PagesController) Update(w http.ResponseWriter, r *http.Request) {
	var page models.Page
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := r.PathValue("page")
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	err = c.app.UpdatePage(r.Context(), &commands.UpdatePageRequest{
		Config: *config,
		Page:   page,
		Path:   path,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"path": path})

}
