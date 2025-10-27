package configs

import (
	"net/http"
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
	_ "github.com/yaien/cultural/internal/modules/configs/interface/migrations"
	"github.com/yaien/cultural/internal/modules/configs/interface/repositories"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"github.com/yaien/cultural/internal/shared"
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
		Configs: repositories.NewConfigRepository(mono.MongoDB),
		Cache:   shared.NewCache[*models.Config](time.Hour),
	})

	m.Web = &Web{
		Middlewares: Middlewares{
			WithConfig: middlewares.NewWithConfig(m.App),
		},
	}

	mono.Router.Handle("/", m.Web.WithConfig(mono.WebRouter))

	return nil

}
