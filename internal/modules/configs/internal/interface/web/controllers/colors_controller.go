package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type ColorsController struct {
	app *application.Application
}

func NewColorsController(app *application.Application) *ColorsController {
	return &ColorsController{app: app}
}

func (c *ColorsController) Get(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	WriteJSON(w, config.Colors)
}

func (c *ColorsController) Update(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	var colors map[string]string
	if err := json.NewDecoder(r.Body).Decode(&colors); err != nil {
		WriteJSONErr(w, models.DecodeError(err))
		return
	}

	if err := c.app.UpdateColors(r.Context(), config, colors); err != nil {
		WriteJSONErr(w, err)
		return
	}

	WriteJSON(w, map[string]any{"success": true})
}
