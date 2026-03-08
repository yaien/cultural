package web

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/routes"
)

type Web struct {
	*middlewares.Middlewares
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
	}

	routes.Register(mono, app, web.Middlewares)
	return web
}
