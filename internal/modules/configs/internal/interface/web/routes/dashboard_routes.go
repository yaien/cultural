package routes

import (
	"net/http"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/assets"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/controllers"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

func dashboard(mono *infrastructure.Monolith, app *application.Application, md *middlewares.Middlewares, reg *models.IntegrationRegistry) {

	{
		ctrl := controllers.NewDashboardController(app)
		mono.WebRouter.Handle("GET /dashboard", md.WithPath(md.WithRole(http.HandlerFunc(ctrl.Home))))
		mono.DashboardRouter.Handle("GET /dashboard/", http.RedirectHandler("/dashboard", http.StatusPermanentRedirect))
		mono.DashboardRouter.HandleFunc("GET /dashboard/empty", md.WithCache(http.HandlerFunc(ctrl.Empty)))
	}

	{
		ctrl := controllers.NewPagesController(app)
		mono.DashboardRouter.HandleFunc("GET /dashboard/pages", ctrl.Index)
		mono.DashboardRouter.HandleFunc("GET /dashboard/pages/preview", ctrl.Preview)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/pages/basic", ctrl.UpdateBasic)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/pages/source", ctrl.UpdateSource)
		mono.DashboardRouter.HandleFunc("POST /dashboard/pages", ctrl.Create)
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/pages", ctrl.Delete)
		mono.DashboardRouter.HandleFunc("POST /dashboard/draft/commit", ctrl.CommitDraft)
	}

	{
		ctrl := controllers.NewFilesController(app)
		mono.DashboardRouter.HandleFunc("POST /dashboard/files", ctrl.Upload)
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/files/{filename}", ctrl.Delete)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/files/{filename}", ctrl.Rename)
		mono.DashboardRouter.HandleFunc("GET /dashboard/files/{filename}", ctrl.Download)

		mono.WebRouter.HandleFunc("GET /assets/dynamic/files/{filename}", ctrl.Download)
	}

	{
		ctrl := controllers.NewFontsController(app)
		mono.DashboardRouter.HandleFunc("GET /dashboard/fonts", ctrl.List)
		mono.DashboardRouter.HandleFunc("POST /dashboard/fonts", ctrl.Update)
	}

	{
		ctrl := controllers.NewColorsController(app)
		mono.DashboardRouter.HandleFunc("POST /dashboard/colors", ctrl.Create)
		mono.DashboardRouter.HandleFunc("PUT /dashboard/colors/{id}", ctrl.Update)
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/colors/{id}", ctrl.Delete)
	}

	{
		ctrl := controllers.NewEventsController(app)
		mono.DashboardRouter.HandleFunc("GET /dashboard/events", ctrl.Index)
	}

	{
		ctrl := controllers.NewRolesController(app, mono.SessionStore)
		mono.DashboardRouter.HandleFunc("GET /dashboard/roles", ctrl.Index)
		mono.DashboardRouter.HandleFunc("POST /dashboard/roles", ctrl.Create)
		mono.DashboardRouter.Handle("GET /dashboard/roles/create", md.WithCache(http.HandlerFunc(ctrl.ShowCreate)))
		mono.DashboardRouter.Handle("GET /dashboard/roles/delete/{id}", md.WithCache(http.HandlerFunc(ctrl.ShowDelete)))
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/roles/{id}", ctrl.Delete)
	}

	{
		ctrl := controllers.NewIntegrationController(app, reg)
		mono.DashboardRouter.HandleFunc("GET /dashboard/integrations", ctrl.Index)
		mono.DashboardRouter.HandleFunc("GET /dashboard/integrations/{integration}", ctrl.Integration)
		mono.DashboardRouter.HandleFunc("GET /dashboard/integrations/{integration}/oauth/connect", ctrl.OAuthLogin)
		mono.DashboardRouter.HandleFunc("GET /dashboard/integrations/{integration}/oauth/callback", ctrl.OauthCallback)
	}

	{
		ctrl := controllers.NewProductsController(app)
		mono.DashboardRouter.HandleFunc("GET /dashboard/products", ctrl.Index)
	}

	mono.WebRouter.HandleFunc("/dashboard/", md.WithPath(md.WithRole(mono.DashboardRouter)))
	mono.WebRouter.Handle("GET /assets/static/dashboard/", http.StripPrefix("/assets/static/dashboard/", md.WithCache(http.FileServer(http.FS(assets.FS)))))

}
