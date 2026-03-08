package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/pages"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ColorsController struct {
	app *application.Application
}

func NewColorsController(app *application.Application) *ColorsController {
	return &ColorsController{app: app}
}

func (c *ColorsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, models.DecodeError(err))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	color, err := c.app.CreateDraftColor(ctx, config.ID)
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("failed to create draft color: %w", err)))
		return
	}

	pages.Color(color).Render(ctx, w)

}

func (c *ColorsController) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, models.DecodeError(err))
		return
	}
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("invalid previousKey: %w", err)))
		return
	}

	req := &commands.UpdateDraftColorRequest{
		ConfigID: config.ID,
		ID:       id,
		Tag:      r.PostForm.Get("tag"),
		Value:    r.PostForm.Get("value"),
	}

	if err := c.app.UpdateDraftColor(ctx, req); err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("failed to update draft colors: %w", err)))
		return
	}

	w.Header().Set("HX-Trigger", "render")
	w.WriteHeader(http.StatusOK)
}

func (c *ColorsController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("invalid id: %w", err)))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	err = c.app.DeleteDraftColor(ctx, config.ID, id)
	if err != nil {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("failed to delete draft color: %w", err)))
		return
	}

	w.Header().Set("HX-Trigger", "render")
	w.WriteHeader(http.StatusOK)
}
