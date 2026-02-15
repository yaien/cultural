package routes

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/controllers"
)

func auth(mono *infrastructure.Monolith, app *application.Application) {
	ctrl := controllers.NewAuthController(app, mono.SessionStore)
	mono.WebRouter.HandleFunc("GET /auth/{provider}/login", ctrl.Login)
	mono.WebRouter.HandleFunc("GET /auth/{provider}/callback", ctrl.Callback)
	mono.WebRouter.HandleFunc("GET /auth/logout", ctrl.Logout)
}
