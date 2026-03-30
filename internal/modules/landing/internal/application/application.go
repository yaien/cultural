package application

import (
	"github.com/yaien/cultural/internal/cache"
	"github.com/yaien/cultural/internal/modules/landing/internal/application/queries"
)

type Application struct {
	*queries.GetPageTemplateQuery
}

type Deps struct {
	Cache *cache.Cache[string]
}

func New(deps Deps) *Application {
	return &Application{
		GetPageTemplateQuery: queries.NewGetPageTemplateQuery(deps.Cache),
	}
}
