package application

import (
	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/library/mail"
	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type Application struct {
	*queries.GetConfigByHostQuery
	*queries.GetUserByIDQuery
	*queries.GetRoleQuery
	*queries.GetRolesQuery
	*queries.GetFontsQuery
	*queries.GetFileQuery
	*queries.GetFilesQuery

	*commands.CreateInvitationCommand
	*commands.SyncUserCommand
	*commands.AcceptInvitationCommand
	*commands.UpdateRoleCommand
	*commands.DeleteRoleCommand
	*commands.UpdatePageCommand
	*commands.CreatePageCommand
	*commands.DeletePageCommand
	*commands.UpdateFontsCommand
	*commands.UploadFileCommand
	*commands.RenameFileCommand
	*commands.DeleteFileCommand
}

type Deps struct {
	Configs       models.ConfigRepository
	Invitations   models.InvitationRepository
	Organizations models.OrganizationRepository
	Roles         models.RoleRepository
	Groups        models.GroupRepository
	Users         models.UserRepository
	Fonts         models.FontRepository
	Files         models.FileRepository
	Cache         *cache.Cache[*models.Config]
	Mail          mail.Mail
	Storage       storage.Storage
}

func New(deps Deps) *Application {
	return &Application{
		queries.NewGetConfigByHostQuery(deps.Configs, deps.Cache),
		queries.NewGetUserByIDQuery(deps.Users),
		queries.NewGetRoleQuery(deps.Roles),
		queries.NewGetRolesQuery(deps.Roles),
		queries.NewGetFontsQuery(deps.Fonts),
		queries.NewGetFileQuery(deps.Files, deps.Storage),
		queries.NewGetFilesQuery(deps.Files),

		commands.NewCreateInvitationCommand(deps.Invitations, deps.Organizations, deps.Configs, deps.Roles, deps.Groups, deps.Mail),
		commands.NewSyncUserCommand(deps.Users),
		commands.NewAcceptInvitationCommand(deps.Invitations, deps.Roles),
		commands.NewUpdateRoleCommand(deps.Roles, deps.Groups),
		commands.NewDeleteRoleCommand(deps.Roles),
		commands.NewUpdatePageCommand(deps.Configs, deps.Cache),
		commands.NewCreatePageCommand(deps.Configs, deps.Cache),
		commands.NewDeletePageCommand(deps.Configs, deps.Cache),
		commands.NewUpdateFontsCommand(deps.Configs, deps.Cache),
		commands.NewUploadFileCommand(deps.Files, deps.Storage),
		commands.NewRenameFileCommand(deps.Files),
		commands.NewDeleteFileCommand(deps.Files, deps.Storage),
	}
}
