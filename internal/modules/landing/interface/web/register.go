package web

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/controllers"
)

func Register(mono *infrastructure.Monolith) {
	ctrl := controllers.NewIndexController()

	mono.Router.HandleFunc("GET /", ctrl.Index)
	mono.Router.HandleFunc("GET /about", ctrl.About)
}
