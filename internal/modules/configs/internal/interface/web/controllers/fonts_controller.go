package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
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
		WriteJSONErr(w, fmt.Errorf("failed getting fonts: %w", err))
		return
	}

	WriteJSON(w, fonts)

}
