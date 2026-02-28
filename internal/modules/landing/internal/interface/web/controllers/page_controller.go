package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/yaien/cultural/internal/modules/configs"
	"github.com/yaien/cultural/internal/modules/landing/internal/application"
	"github.com/yaien/cultural/internal/modules/landing/internal/interface/web/assets"
)

type PageController struct {
	app *application.Application
}

func NewPageController(app *application.Application) *PageController {
	return &PageController{
		app: app,
	}
}

func (c *PageController) Page(w http.ResponseWriter, r *http.Request) {

	config := r.Context().Value(configs.ConfigContextKey).(*configs.Config)

	path := r.PathValue("page")

	html, found, err := c.app.GetPageHTML(config, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "No page template found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(html))
}

func (c *PageController) PageStyles(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(configs.ConfigContextKey).(*configs.Config)

	path := r.PathValue("page")

	if !strings.HasSuffix(path, ".css") {
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

func (c *PageController) LayoutStyles(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(configs.ConfigContextKey).(*configs.Config)

	path := r.PathValue("layout")

	if !strings.HasSuffix(path, ".css") {
		http.NotFound(w, r)
		return
	}

	path = strings.TrimSuffix(path, ".css")

	layout, ok := config.Layouts[path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	_, _ = w.Write([]byte(layout.Styles))
}

func (c *PageController) BaseStyles(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(configs.ConfigContextKey).(*configs.Config)

	w.Header().Set("Content-Type", "text/css")
	err := configs.WritePageBaseStyles(w, config)
	if err != nil {
		http.Error(w, "Failed to generate styles", http.StatusInternalServerError)
		return
	}
}

func (c *PageController) PageScripts(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(configs.ConfigContextKey).(*configs.Config)

	path := r.PathValue("page")

	if !strings.HasSuffix(path, ".js") {
		http.NotFound(w, r)
		return
	}

	path = strings.TrimSuffix(path, ".js")

	page, ok := config.Pages[path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	script := fmt.Sprintf("(() => {\n%s\n})()", page.Script)

	w.Header().Set("Content-Type", "application/javascript")

	_, _ = w.Write([]byte(script))
}

func (c *PageController) LayoutScripts(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(configs.ConfigContextKey).(*configs.Config)

	path := r.PathValue("layout")

	if !strings.HasSuffix(path, ".js") {
		http.NotFound(w, r)
		return
	}

	path = strings.TrimSuffix(path, ".js")

	layout, ok := config.Layouts[path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	script := fmt.Sprintf("(() => {\n%s\n})()", layout.Script)

	w.Header().Set("Content-Type", "application/javascript")

	_, _ = w.Write([]byte(script))
}

func (c *PageController) Favicon(w http.ResponseWriter, r *http.Request) {
	icon, err := assets.FS.Open("favicon.png")
	if err != nil {
		http.Error(w, "Failed to open favicon", http.StatusInternalServerError)
		return
	}

	defer func() {
		err = icon.Close()
		if err != nil {
			http.Error(w, "Failed to close favicon", http.StatusInternalServerError)
			return
		}
	}()

	stat, err := icon.Stat()
	if err != nil {
		http.Error(w, "Failed to get favicon stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	_, err = io.Copy(w, icon)
	if err != nil {
		http.Error(w, "Failed to copy favicon", http.StatusInternalServerError)
		return
	}

}
