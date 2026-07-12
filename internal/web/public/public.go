package public

import (
	"github.com/yaien/cultural/internal/application"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/web/public/routes"
)

func Register(mono *infrastructure.Monolith, app *application.Application) {
	routes.Register(mono, app)
}
