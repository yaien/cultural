package views

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
)

//go:embed *.html icons/*.svg
var fs embed.FS

var Welcome = compile("welcome.html", "icons/*.svg")
var Home = compile("dashboard.html", "home.html", "icons/*.svg")
var Members = compile("dashboard.html", "members.html", "icons/*.svg")
var Pages = compile("dashboard.html", "pages.html", "icons/*.svg")
var Products = compile("dashboard.html", "products.html", "icons/*.svg")
var Events = compile("dashboard.html", "events.html", "icons/*.svg")

func compile(patterns ...string) func(w http.ResponseWriter, r *http.Request, options ...option) {
	t := template.Must(template.ParseFS(fs, patterns...))
	return func(w http.ResponseWriter, r *http.Request, options ...option) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		err := t.ExecuteTemplate(w, patterns[0], newData(r, options...))
		if err != nil {
			slog.Error("failed to render template", "error", err)
		}
	}
}
