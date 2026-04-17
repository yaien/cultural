package routes

import (
	"github.com/yaien/cultural/internal/application"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/web/dashboard/controllers"
)

func external(mono *infrastructure.Monolith, app *application.Application) {
	ctrl := controllers.NewExternalController(app.Storage)
	mono.Router.HandleFunc("/assets/external/{organization_id}/{filename}", ctrl.GetFile)
}
