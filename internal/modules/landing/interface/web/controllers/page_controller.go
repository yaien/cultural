package controllers

import (
	"net/http"
	"text/template"

	"github.com/yaien/cultural/internal/modules/configs/library/views"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type PageController struct{}

func NewPageController() *PageController {
	return &PageController{}
}

func (c *PageController) Index(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	page, ok := config.Pages["index"]
	if !ok {
		http.Error(w, "No index page configured", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_ = views.Page(page, nil).Render(r.Context(), w)
}

func (c *PageController) Page(w http.ResponseWriter, r *http.Request) {

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	path := r.PathValue("page")

	if path == "index" {
		http.NotFound(w, r)
		return
	}

	page, ok := config.Pages[path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_ = views.Page(page, nil).Render(r.Context(), w)
}

var styles = template.Must(template.New("styles").Parse(`
	:root {
	{{range $key, $value := .Fonts.Families}}
		--font-{{ $key }}: '{{ $value }}', sans-serif;
	{{ end }}		
	{{range $key, $value := .Colors}}
		--color-{{ $key }}: {{ $value }};
	{{ end }}
	}
`))

func (c *PageController) Styles(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	w.Header().Set("Content-Type", "text/css")
	err := styles.Execute(w, config)
	if err != nil {
		http.Error(w, "Failed to generate styles", http.StatusInternalServerError)
		return
	}
}
