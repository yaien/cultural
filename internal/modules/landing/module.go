package landing

import (
	"html/template"
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/landing/application"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/routes"
)

type Module struct {
	App *application.Application
}

func (m *Module) Init(mono *infrastructure.Monolith) error {
	m.App = application.New(application.Deps{
		Cache: cache.New[*template.Template](30 * time.Minute),
	})

	routes.Register(mono, m.App)
	return nil

}
