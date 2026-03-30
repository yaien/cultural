package web

import (
	"github.com/yaien/cultural/internal/application/admin"
	"github.com/yaien/cultural/internal/application/auth"
	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/application/preview"
	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/application/store"
	"github.com/yaien/cultural/internal/lib/cache"
	"github.com/yaien/cultural/internal/lib/mail"
	"github.com/yaien/cultural/internal/lib/worker"
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
