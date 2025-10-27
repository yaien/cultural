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
	mono.WebRouter.HandleFunc("GET /", http.HandlerFunc(ctrl.Index))
	mono.WebRouter.HandleFunc("GET /about", ctrl.About)

}
