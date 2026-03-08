package controllers

import (
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/dashboard"
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

	pages.FontList(fonts, family, limit, offset+limit).Render(r.Context(), w)

}

func (c *FontsController) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	req := commands.UpdateDraftFontRequest{
		ConfigID: config.ID,
		Family:   r.PostForm.Get("family"),
		Tag:      r.PostForm.Get("tag"),
	}

	if err := c.app.UpdateDraftFont(ctx, req); err != nil {
		WriteHTMLErr(w, err)
	}

	w.Header().Set("HX-Trigger", "updated")
	dashboard.Toast("Fuente actualizada correctamente", dashboard.Primary).Render(ctx, w)

}
