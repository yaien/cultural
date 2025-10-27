package web

import (
	"net/http"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/assets"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/controllers"
)

func Register(mono *infrastructure.Monolith) {
	ctrl := controllers.NewIndexController()

	mono.Router.Handle("GET /static/landing/", http.StripPrefix("/static/landing/", http.FileServer(http.FS(assets.FS))))
	mono.Router.HandleFunc("GET /", ctrl.Index)
	mono.Router.HandleFunc("GET /about", ctrl.About)
}
