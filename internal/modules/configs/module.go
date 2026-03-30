package configs

import (
	"time"

	"github.com/yaien/cultural/internal/admin"
	"github.com/yaien/cultural/internal/auth"
	"github.com/yaien/cultural/internal/cache"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/integration"
	"github.com/yaien/cultural/internal/integration/integrations/instagram"
	"github.com/yaien/cultural/internal/label"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/preview"
	"github.com/yaien/cultural/internal/storage"
	"github.com/yaien/cultural/internal/store"
	"github.com/yaien/cultural/internal/worker"

	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web"
)

type Module struct {
	App *application.Application
	Web *web.Web
}

func (m *Module) Init(mono *infrastructure.Monolith) error {

	var deps application.Deps

	deps.Mail = mono.Mail
	deps.Cache = cache.New[*label.Config](time.Hour)
	deps.Queue = mono.Queue
	deps.Storage = storage.New(mono.StorageDriver, storage.NewMongo(mono.MongoDB), mono.Queue)
	deps.Store = store.New(store.NewMongo(mono.MongoDB), deps.Storage)
	deps.Auth = auth.New(auth.NewMongo(mono.MongoDB))
	deps.Admin = admin.New(
		admin.NewMongoRoles(mono.MongoDB),
		admin.NewMongoOrganizations(mono.MongoDB),
		admin.NewMongoInvitations(mono.MongoDB),
		admin.NewMongoGroups(mono.MongoDB),
		mono.Mail,
	)
	deps.Label = label.New(
		label.NewMongoFonts(mono.MongoDB),
		label.NewMongoConfigs(mono.MongoDB),
		label.NewMongoDrafts(mono.MongoDB),
		deps.Cache,
	)
	deps.Registry = integration.NewRegistry(
		instagram.Mew(mono.MongoDB),
	)
	deps.Preview = preview.New(deps.Registry)

	m.App = application.New(deps)

	m.Web = web.Register(mono, m.App)

	mono.Worker.Register(worker.H{
		Name:    storage.TaskName,
		Handler: storage.NewHandler(deps.Storage),
	})

	for _, itgr := range deps.Registry.All() {
		if background, ok := itgr.(integration.Background); ok {
			background.RegisterBackgroundProcess(mono.Cron, mono.Queue, mono.Worker)
		}
	}

	mono.Router.Handle("/", m.Web.WithConfig(mono.WebRouter))

	return nil

}
