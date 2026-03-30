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

	app.Storage = storage.New(mono.StorageDriver, storage.NewMongo(mono.MongoDB), mono.Queue)
	app.Store = store.New(store.NewMongo(mono.MongoDB), app.Storage)
	app.Auth = auth.New(auth.NewMongo(mono.MongoDB))
	app.Admin = admin.New(
		admin.NewMongoRoles(mono.MongoDB),
		admin.NewMongoOrganizations(mono.MongoDB),
		admin.NewMongoInvitations(mono.MongoDB),
		admin.NewMongoGroups(mono.MongoDB),
		mono.Mail,
	)
	app.Label = label.New(
		label.NewMongoFonts(mono.MongoDB),
		label.NewMongoConfigs(mono.MongoDB),
		label.NewMongoDrafts(mono.MongoDB),
		cache.New[*label.Config](time.Hour),
	)
	app.Registry = integration.NewRegistry(
		instagram.Mew(mono.MongoDB),
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
