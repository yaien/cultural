package application

import (
	"github.com/yaien/cultural/internal/modules/configs/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"github.com/yaien/cultural/internal/shared"
)

type Application struct {
	*queries.GetConfigByHostQuery
}

type Deps struct {
	Configs models.ConfigRepostory
	Cache   *shared.Cache[*models.Config]
}

func New(deps Deps) *Application {
	return &Application{
		GetConfigByHostQuery: queries.NewGetConfigByHostQuery(deps.Configs, deps.Cache),
	}
}
