package routes

import (
	"net/http"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/assets"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/controllers"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
)

func dashboard(mono *infrastructure.Monolith, app *application.Application, md *middlewares.Middlewares) {

	router := http.NewServeMux()

	{
		ctrl := controllers.NewDashboardController(app)
		mono.WebRouter.Handle("GET /dashboard", md.WithPath(md.WithRole(http.HandlerFunc(ctrl.Home))))
		router.Handle("GET /dashboard/", http.RedirectHandler("/dashboard", http.StatusPermanentRedirect))
		router.HandleFunc("GET /dashboard/empty", md.WithCache(http.HandlerFunc(ctrl.Empty)))
	}

	{
		ctrl := controllers.NewPagesController(app)
		router.HandleFunc("GET /dashboard/pages", ctrl.Index)
		router.HandleFunc("GET /dashboard/pages/preview", ctrl.Preview)
		router.HandleFunc("PATCH /dashboard/pages/basic", ctrl.UpdateBasic)
		router.HandleFunc("PATCH /dashboard/pages/source", ctrl.UpdateSource)
		router.HandleFunc("POST /dashboard/pages", ctrl.Create)
		router.HandleFunc("DELETE /dashboard/pages", ctrl.Delete)
		router.HandleFunc("POST /dashboard/draft/commit", ctrl.CommitDraft)
	}

	{
		ctrl := controllers.NewFilesController(app)
		router.HandleFunc("POST /dashboard/files", ctrl.Upload)
		router.HandleFunc("DELETE /dashboard/files/{filename}", ctrl.Delete)
		router.HandleFunc("PATCH /dashboard/files/{filename}", ctrl.Rename)

		mono.WebRouter.HandleFunc("GET /assets/dynamic/files/{filename}", ctrl.Download)
	}

	{
		ctrl := controllers.NewFontsController(app)
		router.HandleFunc("GET /dashboard/fonts", ctrl.List)
		router.HandleFunc("POST /dashboard/fonts", ctrl.Update)
	}

	{
		ctrl := controllers.NewColorsController(app)
		router.HandleFunc("POST /dashboard/colors", ctrl.Create)
		router.HandleFunc("PUT /dashboard/colors/{id}", ctrl.Update)
		router.HandleFunc("DELETE /dashboard/colors/{id}", ctrl.Delete)
	}

	{
		ctrl := controllers.NewEventsController(app)
		router.HandleFunc("GET /dashboard/events", ctrl.Index)
	}

	{
		ctrl := controllers.NewRolesController(app, mono.SessionStore)
		router.HandleFunc("GET /dashboard/roles", ctrl.Index)
		router.HandleFunc("POST /dashboard/roles", ctrl.Create)
		router.Handle("GET /dashboard/roles/create", md.WithCache(http.HandlerFunc(ctrl.ShowCreate)))
		router.Handle("GET /dashboard/roles/delete/{id}", md.WithCache(http.HandlerFunc(ctrl.ShowDelete)))
		router.HandleFunc("DELETE /dashboard/roles/{id}", ctrl.Delete)
	}

	{
		ctrl := controllers.NewProductsController(app)
		router.HandleFunc("GET /dashboard/products", ctrl.Index)
	}

	mono.WebRouter.HandleFunc("/dashboard/", md.WithPath(md.WithRole(router)))
	mono.WebRouter.Handle("GET /assets/static/dashboard/", http.StripPrefix("/assets/static/dashboard/", md.WithCache(http.FileServer(http.FS(assets.FS)))))

}
