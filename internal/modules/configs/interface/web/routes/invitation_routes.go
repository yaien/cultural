package routes

import (
	"net/http"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/controllers"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/middlewares"
)

func invitations(mono *infrastructure.Monolith, app *application.Application, md *middlewares.Middlewares) {
	ctrl := controllers.NewInvitationController(app)
	mono.WebRouter.HandleFunc("GET /invitation/{id}", md.WithUser(http.HandlerFunc(ctrl.Accept)))
}
