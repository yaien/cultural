package controllers

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(config.Pages)
}

func (c *PagesController) Create(w http.ResponseWriter, r *http.Request) {

	var page models.Page
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	err = c.app.CreatePage(r.Context(), config, &page)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"path": page.Name})
}

func (c *PagesController) Update(w http.ResponseWriter, r *http.Request) {
	var page models.Page
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	path := r.PathValue("page")
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	err = c.app.UpdatePage(r.Context(), &commands.UpdatePageRequest{
		Page:   &page,
		Config: config,
		Path:   path,
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"path": path})

}

func (c *PagesController) Delete(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("page")
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	err := c.app.DeletePage(r.Context(), config, path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"path": path})
}
