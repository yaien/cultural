package routes

import (
	"net/http"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/assets"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/controllers"
)

func Register(mono *infrastructure.Monolith) {
	ctrl := controllers.NewPageController()

	mono.WebRouter.Handle("GET /assets/static/landing/", http.StripPrefix("/assets/static/landing/", http.FileServer(http.FS(assets.FS))))
	mono.WebRouter.HandleFunc("GET /assets/landing/styles.css", ctrl.Styles)
	mono.WebRouter.HandleFunc("GET /{page}", ctrl.Page)
	mono.WebRouter.HandleFunc("GET /", ctrl.Index)
}
