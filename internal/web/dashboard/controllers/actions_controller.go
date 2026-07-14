package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/yaien/cultural/internal/application/label"
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
		ConfigID: config.ID,
		PageName: r.PostForm.Get(pages.SelectedKeyQuery),
		Function: r.PostForm.Get("function"),
	}

	if err := c.drafts.CreateAction(ctx, req); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed creating action: %w", err))
		return
	}

	query := url.Values{}
	query.Set(pages.SelectedKeyQuery, req.PageName)
	query.Set(pages.SectionQuery, pages.EditActionSection)
	query.Set("function", req.Function)

	w.Header().Set("HX-Redirect", fmt.Sprintf("%s?%s", pages.Path, query.Encode()))
	w.WriteHeader(http.StatusCreated)
}

func (c *ActionsController) Update(w http.ResponseWriter, r *http.Request) {

}

func (c *ActionsController) Delete(w http.ResponseWriter, r *http.Request) {

}
