package configs

import (
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/library/worker"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/integrations/instagram"
	_ "github.com/yaien/cultural/internal/modules/configs/internal/interface/migrations"

	"github.com/yaien/cultural/internal/modules/configs/internal/interface/repositories"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type Module struct {
	App      *application.Application
	Web      *web.Web
	Registry *models.IntegrationRegistry
}

func (m *Module) Init(mono *infrastructure.Monolith) error {
	m.Registry = models.NewIntegrationRegistry(
		instagram.Mew(mono.MongoDB),
	)

	deps := application.Deps{
		Configs:       repositories.NewConfigRepository(mono.MongoDB),
		Invitations:   repositories.NewInvitationRepository(mono.MongoDB),
		Organizations: repositories.NewOrganizationRepository(mono.MongoDB),
		Roles:         repositories.NewRoleRepository(mono.MongoDB),
		Groups:        repositories.NewGroupRepository(mono.MongoDB),
		Users:         repositories.NewUserRepository(mono.MongoDB),
		Fonts:         repositories.NewFontRepository(mono.MongoDB),
		Drafts:        repositories.NewDraftRepository(mono.MongoDB),
		Products:      repositories.NewProductRepository(mono.MongoDB),
		Registry:      m.Registry,
		Cache:         cache.New[*models.Config](time.Hour),
		Mail:          mono.Mail,
		Storage:       storage.New(mono.StorageDriver, storage.NewMongo(mono.MongoDB), mono.Queue),
		Queue:         mono.Queue,
	}

	m.App = application.New(deps)

	m.Web = web.Register(mono, m.App, m.Registry)

	mono.Worker.Register(worker.H{
		Name:    storage.TaskName,
		Handler: storage.NewHandler(deps.Storage),
	})

	for _, integration := range m.Registry.All() {
		if background, ok := integration.(models.IntegrationBackground); ok {
			background.RegisterBackgroundProcess(mono.Cron, mono.Queue, mono.Worker)
		}
	}

	mono.Router.Handle("/", m.Web.WithConfig(mono.WebRouter))

	return nil

}
