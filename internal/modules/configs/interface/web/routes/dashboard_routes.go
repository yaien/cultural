package routes

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/assets"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/controllers"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/views"
)

func dashboard(mono *infrastructure.Monolith, app *application.Application, md *middlewares.Middlewares) {
	ctrl := controllers.NewDashboardController(app)

	router := http.NewServeMux()
	router.Handle("GET /dashboard/sites", templ.Handler(views.Sites()))
	router.Handle("GET /dashboard/events", templ.Handler(views.Events()))
	router.Handle("GET /dashboard/products", templ.Handler(views.Products()))
	router.Handle("GET /dashboard/members", templ.Handler(views.Members()))

	mono.WebRouter.Handle("GET /assets/static/dashboard/", http.StripPrefix("/assets/static/dashboard/", http.FileServer(http.FS(assets.FS))))
	mono.WebRouter.HandleFunc("GET /dashboard", md.WithUser(md.WithRole(http.HandlerFunc(ctrl.Home))))
	mono.WebRouter.HandleFunc("GET /dashboard/", md.WithUser(md.WithRole(router)))

}
