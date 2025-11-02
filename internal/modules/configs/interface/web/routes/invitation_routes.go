package routes

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/controllers"
)

func invitations(mono *infrastructure.Monolith, app *application.Application) {
	ctrl := controllers.NewInvitationController(app)
	mono.WebRouter.HandleFunc("GET /invitation/{id}", ctrl.OnInvitation)
}
