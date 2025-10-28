package web

import (
	"net/http"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/assets"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/controllers"
)

func Register(mono *infrastructure.Monolith) {
	ctrl := controllers.NewIndexController()

	mono.WebRouter.Handle("GET /static/landing/", http.StripPrefix("/static/landing/", http.FileServer(http.FS(assets.FS))))
	mono.WebRouter.HandleFunc("GET /{site}", ctrl.Site)
	mono.WebRouter.HandleFunc("GET /", ctrl.Index)
}
