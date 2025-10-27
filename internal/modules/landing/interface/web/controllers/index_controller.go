package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/landing/interface/web/views"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (c *IndexController) Index(w http.ResponseWriter, r *http.Request) {
	// Set content type to HTML
	w.Header().Set("Content-Type", "text/html")

	// Render the index template
	component := views.Index()
	_ = component.Render(r.Context(), w)
}

func (c *IndexController) About(w http.ResponseWriter, r *http.Request) {
	// Set content type to HTML
	w.Header().Set("Content-Type", "text/html")

	// Render the about template
	component := views.About()
	_ = component.Render(r.Context(), w)
}
