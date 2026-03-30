package web

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/routes"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type Web struct {
	*middlewares.Middlewares
}

func Register(mono *infrastructure.Monolith, app *application.Application, registry *models.IntegrationRegistry) *Web {
	web := &Web{
		Middlewares: &middlewares.Middlewares{
			WithConfig: middlewares.NewWithConfig(app),
			WithUser:   middlewares.NewWithUser(app.Deps.Auth.Users, mono.SessionStore),
			WithRole:   middlewares.NewWithRole(app.Deps.Admin.Roles, mono.SessionStore),
			WithCache:  middlewares.WithCache,
			WithPath:   middlewares.WithPath,
		},
	}

	routes.Register(mono, app, web.Middlewares, registry)
	return web
}
