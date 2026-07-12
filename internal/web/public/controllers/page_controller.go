package controllers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/web/middlewares"
	"github.com/yaien/cultural/internal/web/public/assets"
)

type PageController struct {
	registry *integration.Registry
}

func NewPageController(rg *integration.Registry) *PageController {
	return &PageController{rg}
}

func (q *PageController) GetPageHTML(ctx context.Context, config *label.Config, path string) (html string, found bool, err error) {

	page, params, found := label.GetPageInMap(config.Pages, path)
	if !found {
		return "", false, nil
	}

	layout, ok := config.Layouts[page.Layout]
	if !ok {
		layout = label.DefaultLayout
	}

	funcs := make(template.FuncMap)
	var preset any

	for _, itg := range q.registry.All() {
		if itg, ok := itg.(integration.Template); ok {
			maps.Copy(funcs, itg.TemplateFuncMap(ctx, config))
			presets := itg.TemplatePresetMap(ctx, config)
			if p, ok := presets[page.Preset]; ok {
				preset, err = p.Load(params...)
				if err != nil {
					return "", false, fmt.Errorf("failed at loading preset: %w", err)
				}
			}
		}
	}

	html, err = label.RenderPage(&label.PageData{
		Page:     page,
		Layout:   layout,
		AppTitle: config.Title,
		Fonts:    config.Fonts,
		Colors:   config.Colors,
		Funcs:    funcs,
		Preset:   preset,
	})

	if err != nil {
		return "", false, fmt.Errorf("failed at rendering page: %w", err)
	}

	return html, true, nil
}

func (c *PageController) Page(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	path := r.PathValue("page")

	html, found, err := c.GetPageHTML(ctx, config, path)
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
	config := r.Context().Value(middlewares.ConfigContextKey).(*label.Config)

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
	config := r.Context().Value(middlewares.ConfigContextKey).(*label.Config)

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
	config := r.Context().Value(middlewares.ConfigContextKey).(*label.Config)

	w.Header().Set("Content-Type", "text/css")
	err := label.WritePageBaseStyles(w, config)
	if err != nil {
		http.Error(w, "Failed to generate styles", http.StatusInternalServerError)
		return
	}
}

func (c *PageController) PageScripts(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(middlewares.ConfigContextKey).(*label.Config)

	path := r.PathValue("page")

	if !strings.HasSuffix(path, ".js") {
		http.NotFound(w, r)
		return
	}

	path = strings.TrimSuffix(path, ".js")

	slog.Debug("page script", "path", path)

	page, ok := config.Pages[path]
	if !ok {
		slog.Error("page not found", "path", path)
		http.NotFound(w, r)
		return
	}

	script := fmt.Sprintf("(() => {\n%s\n})()", page.Script)

	w.Header().Set("Content-Type", "application/javascript")

	_, _ = w.Write([]byte(script))
}

func (c *PageController) LayoutScripts(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(middlewares.ConfigContextKey).(*label.Config)

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
