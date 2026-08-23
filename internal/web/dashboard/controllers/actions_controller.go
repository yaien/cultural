package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/web/dashboard/views/dashboard"
	"github.com/yaien/cultural/internal/web/dashboard/views/pages"
	"github.com/yaien/cultural/internal/web/middlewares"
)

type ActionsController struct {
	drafts *label.Drafts
}

func NewActionsController(drafts *label.Drafts) *ActionsController {
	return &ActionsController{
		drafts: drafts,
	}
}

func (c *ActionsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	req := &label.CreateActionRequest{
		ConfigID:  config.ID,
		ModelKey:  r.PostForm.Get(pages.SelectedKeyQuery),
		ModelType: label.DraftModelType(r.PostForm.Get(pages.SelectedTypeQuery)),
		Function:  r.PostForm.Get("function"),
	}

	if err := c.drafts.CreateAction(ctx, req); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed creating action: %w", err))
		return
	}

	query := url.Values{}
	query.Set(pages.SelectedKeyQuery, req.ModelKey)
	query.Set(pages.SelectedTypeQuery, string(req.ModelType))
	query.Set(pages.SectionQuery, pages.EditActionSection)
	query.Set("function", req.Function)

	w.Header().Set("HX-Redirect", fmt.Sprintf("%s?%s", pages.Path, query.Encode()))
	w.WriteHeader(http.StatusCreated)
}

func (c *ActionsController) Header(w http.ResponseWriter, r *http.Request) {
	_ = pages.ActionHeader("", "").Render(r.Context(), w)
}

func (c *ActionsController) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed parsing form: %w", err))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	headers := make(map[string]string)

	names := r.PostForm["headers_names"]
	values := r.PostForm["headers_values"]

	if len(names) != len(values) {
		WriteHTMLErr(w, fmt.Errorf("header name and value count mismatch"))
		return
	}

	for i, name := range names {
		headers[name] = values[i]
	}

	req := label.UpdateActionRequest{
		ConfigID:       config.ID,
		Headers:        headers,
		ModelType:      label.DraftModelType(r.PostForm.Get(pages.SelectedTypeQuery)),
		ModelKey:       r.PostForm.Get(pages.SelectedKeyQuery),
		TargetFunction: r.PathValue("function"),
		Function:       r.PostForm.Get("function"),
		Body:           r.PostForm.Get("body"),
	}

	if err := c.drafts.UpdateAction(ctx, &req); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed updating action: %w", err))
		return
	}

	_ = dashboard.Toast("Acción guardada correctamente", dashboard.Success).Render(r.Context(), w)

}

func (c *ActionsController) Delete(w http.ResponseWriter, r *http.Request) {

}
