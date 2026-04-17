package routes

import (
	"net/http"

	"github.com/yaien/cultural/internal/application"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/web/dashboard/controllers"
	"github.com/yaien/cultural/internal/web/middlewares"
)

func invitations(mono *infrastructure.Monolith, app *application.Application, md *middlewares.Middlewares) {
	ctrl := controllers.NewInvitationController(app.Admin.Invitations)
	mono.WebRouter.HandleFunc("GET /invitation/{id}", md.WithUser(http.HandlerFunc(ctrl.Accept)))
}
