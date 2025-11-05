package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
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

func (c *PagesController) Render(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Type  string       `json:"type"`
		Body  models.Page  `json:"body"`
		Email models.Email `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var buffer bytes.Buffer

	switch input.Type {
	case "page":
		_ = render.Page(input.Body, nil).Render(r.Context(), &buffer)
	case "email":
		_ = render.Email(input.Email, nil).Render(r.Context(), &buffer)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"html": buffer.String()})

}
