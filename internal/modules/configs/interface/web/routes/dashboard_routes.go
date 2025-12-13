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

	router := http.NewServeMux()

	{
		ctrl := controllers.NewDashboardController(app)
		mono.WebRouter.Handle("GET /dashboard", md.WithUser(md.WithUser(http.HandlerFunc(ctrl.Home))))
	}

	{
		ctrl := controllers.NewFontsController(app)
		router.HandleFunc("GET /dashboard/api/fonts", ctrl.List)
		router.HandleFunc("GET /dashboard/api/fonts/config", ctrl.Get)
		router.HandleFunc("PUT /dashboard/api/fonts/config", ctrl.Update)
	}

	{
		ctrl := controllers.NewRenderController()
		router.HandleFunc("POST /dashboard/api/render", ctrl.Render)
	}

	{
		ctrl := controllers.NewPagesController(app)
		router.HandleFunc("GET /dashboard/pages", ctrl.Index)
		router.HandleFunc("GET /dashboard/api/pages", ctrl.List)
		router.HandleFunc("PUT /dashboard/api/pages/{page}", ctrl.Update)
	}

	{
		ctrl := controllers.NewFileController(app)
		router.HandleFunc("POST /dashboard/api/files", ctrl.Upload)
		router.HandleFunc("GET /dashboard/api/files", ctrl.List)
		router.HandleFunc("GET /dashboard/api/files/{filename}", ctrl.Get)
		router.HandleFunc("DELETE /dashboard/api/files/{filename}", ctrl.Delete)
		router.HandleFunc("GET /dashboard/files/{filename}", ctrl.Download)
	}

	router.Handle("GET /dashboard/events", templ.Handler(views.Events()))
	router.Handle("GET /dashboard/products", templ.Handler(views.Products()))
	router.Handle("GET /dashboard/members", templ.Handler(views.Members()))

	mono.WebRouter.HandleFunc("/dashboard/", md.WithUser(md.WithRole(router)))

	mono.WebRouter.Handle("GET /assets/static/dashboard/", http.StripPrefix("/assets/static/dashboard/", http.FileServer(http.FS(assets.FS))))

}
