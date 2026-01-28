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
		router.HandleFunc("POST /dashboard/api/pages", ctrl.Create)
		router.HandleFunc("PUT /dashboard/api/pages/{page}", ctrl.Update)
		router.HandleFunc("DELETE /dashboard/api/pages/{page}", ctrl.Delete)
	}

	{
		ctrl := controllers.NewFileController(app)
		router.HandleFunc("POST /dashboard/api/files", ctrl.Upload)
		router.HandleFunc("GET /dashboard/api/files", ctrl.List)
		router.HandleFunc("GET /dashboard/api/files/{filename}", ctrl.Get)
		router.HandleFunc("DELETE /dashboard/api/files/{filename}", ctrl.Delete)
		router.HandleFunc("PUT /dashboard/api/files/{filename}", ctrl.Rename)

		mono.WebRouter.HandleFunc("GET /assets/dynamic/files/{filename}", ctrl.Download)
	}

	{
		ctrl := controllers.NewEventsController(app)
		router.HandleFunc("GET /dashboard/events", ctrl.Index)
	}

	{
		ctrl := controllers.NewInvitationController(app)
		router.HandleFunc("POST /dashboard/api/invitations", ctrl.Create)
	}

	{
		ctrl := controllers.NewRolesController(app)
		router.HandleFunc("GET /dashboard/roles", ctrl.Index)
		router.HandleFunc("GET /dashboard/api/roles", ctrl.List)
		router.HandleFunc("PUT /dashboard/api/roles/{id}", ctrl.Update)
		router.HandleFunc("DELETE /dashboard/api/roles/{id}", ctrl.Delete)
	}

	{
		ctrl := controllers.NewProductsController(app)
		router.HandleFunc("GET /dashboard/products", ctrl.Index)
	}

	mono.WebRouter.HandleFunc("/dashboard/", md.WithUser(md.WithRole(router)))

	mono.WebRouter.Handle("GET /assets/static/dashboard/", http.StripPrefix("/assets/static/dashboard/", http.FileServer(http.FS(assets.FS))))

}
