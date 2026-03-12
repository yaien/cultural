package web

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/integrations/instagram"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/routes"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type Web struct {
	*middlewares.Middlewares
	IntegrationRegistry *models.IntegrationRegistry
}

func Register(mono *infrastructure.Monolith, app *application.Application) *Web {
	web := &Web{
		Middlewares: &middlewares.Middlewares{
			WithConfig: middlewares.NewWithConfig(app),
			WithUser:   middlewares.NewWithUser(app, mono.SessionStore),
			WithRole:   middlewares.NewWithRole(app, mono.SessionStore),
			WithCache:  middlewares.WithCache,
			WithPath:   middlewares.WithPath,
		},
		IntegrationRegistry: models.NewIntegrationRegistry(
			instagram.Mew(mono.MongoDB),
		),
	}

	routes.Register(mono, app, web.Middlewares, web.IntegrationRegistry)
	return web
}
