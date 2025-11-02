package routes

import (
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/application"
)

func Register(mono *infrastructure.Monolith, app *application.Application) {
	invitations(mono, app)
}
