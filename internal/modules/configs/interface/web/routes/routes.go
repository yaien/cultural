package routes

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/middlewares"
)

func Register(mono *infrastructure.Monolith, app *application.Application, md *middlewares.Middlewares) {
	auth(mono, app)
	invitations(mono, app, md)
	dashboard(mono, app, md)
}
