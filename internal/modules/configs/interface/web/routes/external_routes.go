package routes

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/controllers"
)

func external(mono *infrastructure.Monolith, app *application.Application) {
	ctrl := controllers.NewExternalController(app)
	mono.Router.HandleFunc("/assets/external/{organization_id}/{filename}", ctrl.GetFile)
}
