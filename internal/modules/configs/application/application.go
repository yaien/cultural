package application

import (
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/library/cache"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"gopkg.in/gomail.v2"
)

type Application struct {
	*queries.GetConfigByHostQuery

	*commands.CreateInvitationCommand
}

type Deps struct {
	Configs     models.ConfigRepostory
	Invitations models.InvitationRepository
	Cache       *cache.Cache[*models.Config]
	Mail        *gomail.Dialer
}

func New(deps Deps) *Application {
	return &Application{
		GetConfigByHostQuery:    queries.NewGetConfigByHostQuery(deps.Configs, deps.Cache),
		CreateInvitationCommand: commands.NewCreateInvitationCommand(deps.Invitations, deps.Mail),
	}
}
