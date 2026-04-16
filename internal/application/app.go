package application

import (
	"time"

	"github.com/yaien/cultural/internal/application/admin"
	"github.com/yaien/cultural/internal/application/auth"
	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/integration/instagram"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/application/preview"
	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/application/store"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/lib/cache"
	"github.com/yaien/cultural/internal/lib/worker"
)

type Application struct {
	Registry *integration.Registry
	Storage  *storage.Storage
	Store    *store.Store
	Admin    *admin.Admin
	Auth     *auth.Auth
	Label    *label.Label
	Preview  *preview.Preview
}

func New(mono *infrastructure.Monolith) *Application {
	var app Application

	app.Storage = storage.New(mono.StorageDriver, storage.NewGorm(mono.GormDB), mono.Queue)
	app.Store = store.New(mono.GormDB, app.Storage)
	app.Auth = auth.New(auth.NewGorm(mono.GormDB))
	app.Admin = admin.New(
		admin.NewGormRoles(mono.GormDB),
		admin.NewGormOrganizations(mono.GormDB),
		admin.NewGormInvitations(mono.GormDB),
		admin.NewGormGroups(mono.GormDB),
		mono.Mail,
	)
	app.Label = label.New(
		label.NewGormFonts(mono.GormDB),
		label.NewGormConfigs(mono.GormDB),
		label.NewGormDrafts(mono.GormDB),
		cache.New[*label.Config](time.Hour),
	)

	app.Registry = integration.NewRegistry(
		instagram.Mew(integration.NewGorm[instagram.Data](mono.GormDB), app.Label.Configs),
	)
	app.Preview = preview.New(app.Registry)

	mono.Worker.Register(worker.H{
		Name:    storage.TaskName,
		Handler: storage.NewHandler(app.Storage),
	})

	for _, itgr := range app.Registry.All() {
		if background, ok := itgr.(integration.Background); ok {
			background.RegisterBackgroundProcess(mono.Cron, mono.Queue, mono.Worker)
		}
	}

	return &app
}
