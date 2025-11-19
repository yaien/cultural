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
	*queries.GetFontsQery

	*commands.CreateInvitationCommand
	*commands.SyncUserCommand
	*commands.AcceptInvitationCommand
	*commands.UpdatePageCommand
	*commands.UpdateFontsCommand
}

type Deps struct {
	Configs       models.ConfigRepository
	Invitations   models.InvitationRepository
	Organizations models.OrganizationRepository
	Roles         models.RoleRepository
	Groups        models.GroupRepository
	Users         models.UserRepository
	Fonts         models.FontRepository
	Cache         *cache.Cache[*models.Config]
	Mail          *gomail.Dialer
}

func New(deps Deps) *Application {
	return &Application{
		queries.NewGetConfigByHostQuery(deps.Configs, deps.Cache),
		queries.NewGetUserByIDQuery(deps.Users),
		queries.NewGetRoleQuery(deps.Roles),
		queries.NewGetFontsQuery(deps.Fonts),

		commands.NewCreateInvitationCommand(deps.Invitations, deps.Organizations, deps.Configs, deps.Roles, deps.Groups, deps.Mail),
		commands.NewSyncUserCommand(deps.Users),
		commands.NewAcceptInvitationCommand(deps.Invitations, deps.Roles),
		commands.NewUpdatePageCommand(deps.Configs, deps.Cache),
		commands.NewUpdateFontsCommand(deps.Configs, deps.Cache),
	}
}
