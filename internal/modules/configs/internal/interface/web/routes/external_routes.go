package routes

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/controllers"
)

func external(mono *infrastructure.Monolith, app *application.Application) {
	ctrl := controllers.NewExternalController(app.Deps.Storage)
	mono.Router.HandleFunc("/assets/external/{organization_id}/{filename}", ctrl.GetFile)
}
