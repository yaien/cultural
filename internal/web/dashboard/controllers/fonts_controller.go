package controllers

import (
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/web/dashboard/views/pages"
	"github.com/yaien/cultural/internal/web/middlewares"
)

type FontsController struct {
	fonts  *label.Fonts
	drafts *label.Drafts
}

func NewFontsController(fonts *label.Fonts, drafts *label.Drafts) *FontsController {
	return &FontsController{fonts, drafts}
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

	fonts, err := c.fonts.Find(r.Context(), &label.FindFontOptions{
		Family: family,
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		http.Error(w, "Failed to retrieve fonts", http.StatusInternalServerError)
		return
	}

	_ = pages.FontList(fonts, family, limit, offset+limit).Render(r.Context(), w)

}

func (c *FontsController) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	req := label.UpdateDraftFontOptions{
		ConfigID: config.ID,
		Family:   r.PostForm.Get("family"),
		Tag:      r.PostForm.Get("tag"),
	}

	if err := c.drafts.UpdateFont(ctx, req); err != nil {
		WriteHTMLErr(w, err)
	}

	w.Header().Set("HX-Trigger", "updated, render")
	w.WriteHeader(http.StatusOK)

}
