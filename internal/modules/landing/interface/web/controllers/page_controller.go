package controllers

import (
	"net/http"
	"strings"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"github.com/yaien/cultural/internal/modules/landing/application"
)

type PageController struct {
	app *application.Application
}

func NewPageController(app *application.Application) *PageController {
	return &PageController{
		app: app,
	}
}

func (c *PageController) Index(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	parsed, found, err := c.app.GetPageTemplate(config, "index")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "No index page configured", http.StatusInternalServerError)
		return
	}

	page := config.Pages["index"]

	data := models.NewPageData(config, page).
		WithFilePath("/assets/dynamic/files/").
		Data()

	w.Header().Set("Content-Type", "text/html")
	parsed.Execute(w, data)

}

func (c *PageController) Page(w http.ResponseWriter, r *http.Request) {

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	path := r.PathValue("page")

	if path == "index" {
		http.NotFound(w, r)
		return
	}

	parsed, found, err := c.app.GetPageTemplate(config, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "No page template found", http.StatusInternalServerError)
		return
	}

	page := config.Pages[path]

	data := models.NewPageData(config, page).
		WithFilePath("/assets/dynamic/landing/").
		Data()

	w.Header().Set("Content-Type", "text/html")
	parsed.Execute(w, data)

}

func (c *PageController) PageStyles(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	path := r.PathValue("page")

	if path == "index" || !strings.HasSuffix(path, ".css") {
		http.NotFound(w, r)
		return
	}

	path = strings.TrimSuffix(path, ".css")

	page, ok := config.Pages[path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	_, _ = w.Write([]byte(page.Styles))
}

func (c *PageController) BaseStyles(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	w.Header().Set("Content-Type", "text/css")
	err := models.PageBaseStyles.Execute(w, config)
	if err != nil {
		http.Error(w, "Failed to generate styles", http.StatusInternalServerError)
		return
	}
}
