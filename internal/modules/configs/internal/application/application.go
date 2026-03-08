package application

import (
	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/library/mail"
	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/library/worker"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type Application struct {
	*queries.GetConfigByHostQuery
	*queries.GetUserByIDQuery
	*queries.GetRoleQuery
	*queries.GetRolesQuery
	*queries.GetFontsQuery
	*queries.GetFontQuery
	*queries.GetFileDataQuery
	*queries.GetFileQuery
	*queries.GetFilesQuery
	*queries.GetDraftByConfigIDQuery
	*queries.GetPreviewQuery

	*commands.CreateInvitationCommand
	*commands.SyncUserCommand
	*commands.AcceptInvitationCommand
	*commands.UpdateRoleCommand
	*commands.DeleteRoleCommand
	*commands.UploadFileCommand
	*commands.RenameFileCommand
	*commands.DeleteFileCommand
	*commands.UpdateDraftBasicCommand
	*commands.UpdateDraftSourceCommand
	*commands.UpdateDraftFontCommand
	*commands.UpdateDraftColorCommand
	*commands.CreateDraftColorCommand
	*commands.DeleteDraftColorCommand
	*commands.CreateDraftModelCommand
	*commands.DeleteDraftModelCommand

	*commands.CommitDraftCommand
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
	Drafts        models.DraftRepository
	Cache         *cache.Cache[*models.Config]
	Queue         *worker.Queue
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
		queries.NewGetFontQuery(deps.Fonts),
		queries.NewGetFileDataQuery(deps.Files, deps.Storage),
		queries.NewGetFileQuery(deps.Files),
		queries.NewGetFilesQuery(deps.Files),
		queries.NewGetDraftByConfigIDQuery(deps.Drafts),
		queries.NewGetPreviewQuery(deps.Drafts),

		commands.NewCreateInvitationCommand(deps.Invitations, deps.Organizations, deps.Configs, deps.Roles, deps.Groups, deps.Mail),
		commands.NewSyncUserCommand(deps.Users),
		commands.NewAcceptInvitationCommand(deps.Invitations, deps.Roles),
		commands.NewUpdateRoleCommand(deps.Roles, deps.Groups),
		commands.NewDeleteRoleCommand(deps.Roles),
		commands.NewUploadFileCommand(deps.Files, deps.Storage, deps.Queue),
		commands.NewRenameFileCommand(deps.Files),
		commands.NewDeleteFileCommand(deps.Files, deps.Storage),
		commands.NewUpdateDraftBasicCommand(deps.Drafts),
		commands.NewUpdateDraftSourceCommand(deps.Drafts),
		commands.NewUpdateDraftFontCommand(deps.Drafts, deps.Fonts),
		commands.NewUpdateDraftColorCommand(deps.Drafts),
		commands.NewCreateDraftColorCommand(deps.Drafts),
		commands.NewDeleteDraftColorCommand(deps.Drafts),
		commands.NewCreateDraftModelCommand(deps.Drafts),
		commands.NewDeleteDraftModelCommand(deps.Drafts),
		commands.NewCommitDraftCommand(deps.Configs, deps.Drafts, deps.Cache),
	}
}
