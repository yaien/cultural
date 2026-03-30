package application

import (
	"github.com/yaien/cultural/internal/admin"
	"github.com/yaien/cultural/internal/auth"
	"github.com/yaien/cultural/internal/cache"
	"github.com/yaien/cultural/internal/integration"
	"github.com/yaien/cultural/internal/label"
	"github.com/yaien/cultural/internal/mail"
	"github.com/yaien/cultural/internal/preview"
	"github.com/yaien/cultural/internal/storage"
	"github.com/yaien/cultural/internal/store"
	"github.com/yaien/cultural/internal/worker"
)

type Application struct {
	Deps Deps
}

type Deps struct {
	Mail     mail.Mail
	Cache    *cache.Cache[*label.Config]
	Queue    *worker.Queue
	Registry *integration.Registry
	Storage  *storage.Storage
	Store    *store.Store
	Admin    *admin.Admin
	Auth     *auth.Auth
	Label    *label.Label
	Preview  *preview.Preview
}

func New(deps Deps) *Application {
	return &Application{
		deps,
	}
}
