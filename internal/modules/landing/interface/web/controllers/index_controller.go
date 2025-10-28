package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/views"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (c *IndexController) Site(w http.ResponseWriter, r *http.Request) {

	config := r.Context().Value(middlewares.ConfigContextKey).(*models.Config)

	path := r.PathValue("site")

	site, ok := config.Sites[path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_ = views.Site(site).Render(r.Context(), w)
}

func (c *IndexController) Index(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(middlewares.ConfigContextKey).(*models.Config)
	w.Header().Set("Content-Type", "text/html")
	_ = views.Site(config.Index).Render(r.Context(), w)
}
