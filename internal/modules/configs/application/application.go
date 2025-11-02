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
	*queries.GetUserByIDQuery
	*queries.GetRoleQuery

	*commands.CreateInvitationCommand
	*commands.SyncUserCommand
}

type Deps struct {
	Configs       models.ConfigRepository
	Invitations   models.InvitationRepository
	Organizations models.OrganizationRepository
	Roles         models.RoleRepository
	Groups        models.GroupRepository
	Users         models.UserRepository
	Cache         *cache.Cache[*models.Config]
	Mail          *gomail.Dialer
}

func New(deps Deps) *Application {
	return &Application{
		queries.NewGetConfigByHostQuery(deps.Configs, deps.Cache),
		queries.NewGetUserByIDQuery(deps.Users),
		queries.NewGetRoleQuery(deps.Roles),

		commands.NewCreateInvitationCommand(deps.Invitations, deps.Organizations, deps.Configs, deps.Roles, deps.Groups, deps.Mail),
		commands.NewSyncUserCommand(deps.Users),
	}
}
