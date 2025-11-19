package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type FontsController struct {
	app *application.Application
}

func NewFontsController(app *application.Application) *FontsController {
	return &FontsController{app}
}

func (c *FontsController) List(w http.ResponseWriter, r *http.Request) {
	options := &models.FindFontOptions{}

	options.Family = r.URL.Query().Get("family")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err == nil {
		options.Limit = int64(limit)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err == nil {
		options.Offset = int64(offset)
	}

	fonts, err := c.app.GetFonts(r.Context(), options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fonts)

}

func (c *FontsController) Get(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(config.Fonts)
}

func (c *FontsController) Update(w http.ResponseWriter, r *http.Request) {
	var fonts models.Fonts
	err := json.NewDecoder(r.Body).Decode(&fonts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	err = c.app.UpdateFonts(r.Context(), *config, fonts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}
