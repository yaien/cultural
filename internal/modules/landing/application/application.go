package application

import (
	"html/template"

	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/landing/application/queries"
)

type Application struct {
	*queries.GetPageTemplateQuery
}

type Deps struct {
	Cache *cache.Cache[*template.Template]
}

func New(deps Deps) *Application {
	return &Application{
		GetPageTemplateQuery: queries.NewGetPageTemplateQuery(deps.Cache),
	}
}
