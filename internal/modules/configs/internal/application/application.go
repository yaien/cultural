package application

import (
	"github.com/yaien/cultural/internal/library/admin"
	"github.com/yaien/cultural/internal/library/auth"
	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/library/mail"
	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/library/store"
	"github.com/yaien/cultural/internal/library/worker"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type Application struct {
	Deps Deps

	*queries.GetConfigByHostQuery
	*queries.GetFontsQuery
	*queries.GetFontQuery
	*queries.GetDraftByConfigIDQuery
	*queries.GetPreviewQuery

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
	Configs  models.ConfigRepository
	Fonts    models.FontRepository
	Drafts   models.DraftRepository
	Cache    *cache.Cache[*models.Config]
	Queue    *worker.Queue
	Registry *models.IntegrationRegistry
	Mail     mail.Mail
	Storage  *storage.Storage
	Store    *store.Store
	Admin    *admin.Admin
	Auth     *auth.Auth
}

func New(deps Deps) *Application {
	return &Application{
		deps,

		queries.NewGetConfigByHostQuery(deps.Configs, deps.Cache),
		queries.NewGetFontsQuery(deps.Fonts),
		queries.NewGetFontQuery(deps.Fonts),
		queries.NewGetDraftByConfigIDQuery(deps.Drafts),
		queries.NewGetPreviewQuery(deps.Drafts, deps.Registry),

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
