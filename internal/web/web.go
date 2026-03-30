package web

import (
	"github.com/yaien/cultural/internal/application"
	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/lib/worker"
	"github.com/yaien/cultural/internal/web/dashboard"
	"github.com/yaien/cultural/internal/web/middlewares"
	"github.com/yaien/cultural/internal/web/public"
)

func Register(mono *infrastructure.Monolith, app *application.Application) error {

	mdl := &middlewares.Middlewares{
		WithConfig: middlewares.NewWithConfig(app.Label.Configs),
		WithUser:   middlewares.NewWithUser(app.Auth.Users, mono.SessionStore),
		WithRole:   middlewares.NewWithRole(app.Admin.Roles, mono.SessionStore),
		WithCache:  middlewares.WithCache,
		WithPath:   middlewares.WithPath,
	}

	dashboard.Register(mono, app, mdl)
	public.Register(mono)

	mono.Worker.Register(worker.H{
		Name:    storage.TaskName,
		Handler: storage.NewHandler(app.Storage),
	})

	for _, itgr := range app.Registry.All() {
		if background, ok := itgr.(integration.Background); ok {
			background.RegisterBackgroundProcess(mono.Cron, mono.Queue, mono.Worker)
		}
	}

	mono.Router.Handle("/", mdl.WithConfig(mono.WebRouter))

	return nil

}
