package routes

import (
	"net/http"

	"github.com/yaien/cultural/internal/application"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/web/dashboard/assets"
	"github.com/yaien/cultural/internal/web/dashboard/controllers"
	"github.com/yaien/cultural/internal/web/middlewares"
)

func dashboard(mono *infrastructure.Monolith, app *application.Application, md *middlewares.Middlewares) {

	{
		ctrl := controllers.NewDashboardController()
		mono.WebRouter.Handle("GET /dashboard", md.WithPath(md.WithRole(http.HandlerFunc(ctrl.Home))))
		mono.DashboardRouter.Handle("GET /dashboard/", http.RedirectHandler("/dashboard", http.StatusPermanentRedirect))
		mono.DashboardRouter.HandleFunc("GET /dashboard/empty", md.WithCache(http.HandlerFunc(ctrl.Empty)))
	}

	{
		ctrl := controllers.NewPagesController(app.Label.Drafts, app.Label.Fonts, app.Preview, app.Storage)
		mono.DashboardRouter.HandleFunc("GET /dashboard/pages", ctrl.Index)
		mono.DashboardRouter.HandleFunc("GET /dashboard/pages/preview", ctrl.Preview)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/pages/basic", ctrl.UpdateBasic)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/pages/source", ctrl.UpdateSource)
		mono.DashboardRouter.HandleFunc("POST /dashboard/pages", ctrl.Create)
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/pages", ctrl.Delete)
		mono.DashboardRouter.HandleFunc("POST /dashboard/draft/commit", ctrl.CommitDraft)
	}

	{
		ctrl := controllers.NewFilesController(app.Storage)
		mono.DashboardRouter.HandleFunc("POST /dashboard/files", ctrl.Upload)
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/files/{filename}", ctrl.Delete)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/files/{filename}", ctrl.Rename)
		mono.DashboardRouter.HandleFunc("GET /dashboard/files/{filename}", ctrl.Download)

		mono.WebRouter.HandleFunc("GET /assets/dynamic/files/{filename}", ctrl.Download)
	}

	{
		ctrl := controllers.NewFontsController(app.Label.Fonts, app.Label.Drafts)
		mono.DashboardRouter.HandleFunc("GET /dashboard/fonts", ctrl.List)
		mono.DashboardRouter.HandleFunc("POST /dashboard/fonts", ctrl.Update)
	}

	{
		ctrl := controllers.NewColorsController(app.Label.Drafts)
		mono.DashboardRouter.HandleFunc("POST /dashboard/colors", ctrl.Create)
		mono.DashboardRouter.HandleFunc("PUT /dashboard/colors/{id}", ctrl.Update)
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/colors/{id}", ctrl.Delete)
	}

	{
		ctrl := controllers.NewEventsController()
		mono.DashboardRouter.HandleFunc("GET /dashboard/events", ctrl.Index)
	}

	{
		ctrl := controllers.NewRolesController(app.Admin.Roles, app.Admin.Invitations, mono.SessionStore)
		mono.DashboardRouter.HandleFunc("GET /dashboard/roles", ctrl.Index)
		mono.DashboardRouter.HandleFunc("POST /dashboard/roles", ctrl.Create)
		mono.DashboardRouter.Handle("GET /dashboard/roles/create", md.WithCache(http.HandlerFunc(ctrl.ShowCreate)))
		mono.DashboardRouter.Handle("GET /dashboard/roles/delete/{id}", md.WithCache(http.HandlerFunc(ctrl.ShowDelete)))
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/roles/{id}", ctrl.Delete)
	}

	{
		ctrl := controllers.NewIntegrationController(app.Registry)
		mono.DashboardRouter.HandleFunc("GET /dashboard/integrations", ctrl.Index)
		mono.DashboardRouter.HandleFunc("GET /dashboard/integrations/{integration}", ctrl.Integration)
		mono.DashboardRouter.HandleFunc("GET /dashboard/integrations/{integration}/oauth/connect", ctrl.OAuthLogin)
		mono.DashboardRouter.HandleFunc("GET /dashboard/integrations/{integration}/oauth/callback", ctrl.OauthCallback)
	}

	{
		ctrl := controllers.NewProductsController(app.Store.Products, app.Store.Presentations, app.Store.Contents)
		mono.DashboardRouter.HandleFunc("GET /dashboard/products", ctrl.Index)
		mono.DashboardRouter.HandleFunc("GET /dashboard/products/create", ctrl.CreateModal)
		mono.DashboardRouter.HandleFunc("POST /dashboard/products", ctrl.Create)
		mono.DashboardRouter.HandleFunc("GET /dashboard/products/{id}", ctrl.Show)
		mono.DashboardRouter.HandleFunc("POST /dashboard/products/{id}/presentations", ctrl.CreatePresentation)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/products/{id}/presentations/{presentationId}", ctrl.UpdatePresentation)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/products/{id}/presentations/toggle", ctrl.TogglePresentations)
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/products/{id}/presentations/{presentationId}", ctrl.DeletePresentation)
		mono.DashboardRouter.HandleFunc("POST /dashboard/products/{id}/presentations/{presentationId}/files", ctrl.UploadPresentationFile)
		mono.DashboardRouter.HandleFunc("PATCH /dashboard/products/{id}/presentations/{presentationId}/files/toggle", ctrl.TogglePresentationFiles)
		mono.DashboardRouter.HandleFunc("DELETE /dashboard/products/{id}/presentations/{presentationId}/files/{contentId}", ctrl.DeletePresentationFile)
		mono.DashboardRouter.HandleFunc("GET /dashboard/products/{id}/presentations/{presentationId}/files/{contentId}", ctrl.PickPresentationFile)
	}

	mono.WebRouter.HandleFunc("/dashboard/", md.WithPath(md.WithRole(mono.DashboardRouter)))
	mono.WebRouter.Handle("GET /assets/static/dashboard/", http.StripPrefix("/assets/static/dashboard/", md.WithCache(http.FileServer(http.FS(assets.FS)))))

}
