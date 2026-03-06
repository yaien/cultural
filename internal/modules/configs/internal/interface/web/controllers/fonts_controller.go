package controllers

import (
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/pages"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type FontsController struct {
	app *application.Application
}

func NewFontsController(app *application.Application) *FontsController {
	return &FontsController{app: app}
}

func (c *FontsController) List(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	family := query.Get("family")

	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
		return
	}

	fonts, err := c.app.GetFonts(r.Context(), &models.FindFontOptions{
		Family: family,
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		http.Error(w, "Failed to retrieve fonts", http.StatusInternalServerError)
		return
	}

	pages.FontList(fonts, limit, offset+limit).Render(r.Context(), w)

}
