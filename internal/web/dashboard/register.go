package dashboard

import (
	"github.com/yaien/cultural/internal/application"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/web/dashboard/routes"
	"github.com/yaien/cultural/internal/web/middlewares"
)

func Register(mono *infrastructure.Monolith, app *application.Application, mdl *middlewares.Middlewares) {
	routes.Register(mono, app, mdl)
}
