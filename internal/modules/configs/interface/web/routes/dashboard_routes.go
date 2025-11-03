package routes

import (
	"net/http"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/assets"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/controllers"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/middlewares"
)

func dashboard(mono *infrastructure.Monolith, app *application.Application, md *middlewares.Middlewares) {
	ctrl := controllers.NewDashboardController(app)

	router := http.NewServeMux()

	mono.WebRouter.Handle("GET /assets/static/dashboard/", http.StripPrefix("/assets/static/dashboard/", http.FileServer(http.FS(assets.FS))))
	mono.WebRouter.HandleFunc("GET /dashboard", md.WithUser(md.WithRole(http.HandlerFunc(ctrl.Home))))
	mono.WebRouter.HandleFunc("GET /dashboard/", md.WithUser(md.WithRole(router)))

}
