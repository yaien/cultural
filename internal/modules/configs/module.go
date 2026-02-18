package configs

import (
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	_ "github.com/yaien/cultural/internal/modules/configs/internal/interface/migrations"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/repositories"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type Module struct {
	App *application.Application
	Web *web.Web
}

func (m *Module) Init(mono *infrastructure.Monolith) error {
	m.App = application.New(application.Deps{
		Configs:       repositories.NewConfigRepository(mono.MongoDB),
		Invitations:   repositories.NewInvitationRepository(mono.MongoDB),
		Organizations: repositories.NewOrganizationRepository(mono.MongoDB),
		Roles:         repositories.NewRoleRepository(mono.MongoDB),
		Groups:        repositories.NewGroupRepository(mono.MongoDB),
		Users:         repositories.NewUserRepository(mono.MongoDB),
		Fonts:         repositories.NewFontRepository(mono.MongoDB),
		Files:         repositories.NewFileRepository(mono.MongoDB),
		Drafts:        repositories.NewDraftRepository(mono.MongoDB),
		Cache:         cache.New[*models.Config](time.Hour),
		Mail:          mono.Mail,
		Storage:       mono.Storage,
		Queue:         mono.Queue,
	})

	m.Web = web.Register(mono, m.App)

	mono.Router.Handle("/", m.Web.WithConfig(mono.WebRouter))

	return nil

}
