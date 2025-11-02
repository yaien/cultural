package configs

import (
	"net/http"
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
	_ "github.com/yaien/cultural/internal/modules/configs/interface/migrations"
	"github.com/yaien/cultural/internal/modules/configs/interface/repositories"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/routes"
	"github.com/yaien/cultural/internal/modules/configs/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type Middlewares struct {
	WithConfig func(next http.Handler) http.HandlerFunc
}

type Web struct {
	Middlewares
}

type Module struct {
	App *application.Application
	Web *Web
}

func (m *Module) Init(mono *infrastructure.Monolith) error {
	m.App = application.New(application.Deps{
		Configs:       repositories.NewConfigRepository(mono.MongoDB),
		Invitations:   repositories.NewInvitationRepository(mono.MongoDB),
		Organizations: repositories.NewOrganizationRepository(mono.MongoDB),
		Roles:         repositories.NewRoleRepository(mono.MongoDB),
		Groups:        repositories.NewGroupRepository(mono.MongoDB),
		Cache:         cache.New[*models.Config](time.Hour),
		Mail:          mono.Mail,
	})

	m.Web = &Web{
		Middlewares: Middlewares{
			WithConfig: middlewares.NewWithConfig(m.App),
		},
	}

	mono.Router.Handle("/", m.Web.WithConfig(mono.WebRouter))

	routes.Register(mono, m.App)

	return nil

}
