package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/web/dashboard/views/pages"
	"github.com/yaien/cultural/internal/web/middlewares"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColorsController struct {
	drafts *label.Drafts
}

func NewColorsController(drafts *label.Drafts) *ColorsController {
	return &ColorsController{drafts: drafts}
}

func (c *ColorsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "failed to parse form: %w", err))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	color, err := c.drafts.CreateColor(ctx, config.ID)
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "failed to create draft color: %w", err))
		return
	}

	_ = pages.Color(color).Render(ctx, w)

}

func (c *ColorsController) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "failed to parse form: %w", err))
		return
	}
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid previousKey: %w", err))
		return
	}

	req := &label.UpdateDraftColorOptions{
		ConfigID: config.ID,
		ID:       id,
		Tag:      r.PostForm.Get("tag"),
		Value:    r.PostForm.Get("value"),
	}

	if err := c.drafts.UpdateColor(ctx, req); err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "failed to update draft colors: %w", err))
		return
	}

	w.Header().Set("HX-Trigger", "render")
	w.WriteHeader(http.StatusOK)
}

func (c *ColorsController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "invalid id: %w", err))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	err = c.drafts.DeleteColor(ctx, config.ID, id)
	if err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "failed to delete draft color: %w", err))
		return
	}

	w.Header().Set("HX-Trigger", "render")
	w.WriteHeader(http.StatusOK)
}
