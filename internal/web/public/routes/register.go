package routes

import (
	"net/http"
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/lib/cache"
	"github.com/yaien/cultural/internal/web/dashboard/assets"
	"github.com/yaien/cultural/internal/web/public/controllers"
)

func Register(mono *infrastructure.Monolith) {
	ctrl := controllers.NewPageController(cache.New[string](30 * time.Minute))

	mono.WebRouter.Handle("GET /assets/static/landing/", http.StripPrefix("/assets/static/landing/", http.FileServer(http.FS(assets.FS))))
	mono.WebRouter.HandleFunc("GET /assets/landing/styles.css", ctrl.BaseStyles)
	mono.WebRouter.HandleFunc("GET /assets/landing/favicon.png", ctrl.Favicon)
	mono.WebRouter.HandleFunc("GET /assets/landing/styles/pages/{page}", ctrl.PageStyles)
	mono.WebRouter.HandleFunc("GET /assets/landing/styles/layouts/{layout}", ctrl.LayoutStyles)
	mono.WebRouter.HandleFunc("GET /assets/landing/scripts/pages/{page}", ctrl.PageScripts)
	mono.WebRouter.HandleFunc("GET /assets/landing/scripts/layouts/{layout}", ctrl.LayoutScripts)
	mono.WebRouter.HandleFunc("/{page...}", ctrl.Page)
	mono.WebRouter.HandleFunc("/{$}", ctrl.Page)
}
