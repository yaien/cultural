package queries

import (
	"fmt"
	"html/template"

	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/configs"
)

type GetPageTemplateQuery struct {
	cache *cache.Cache[*template.Template]
}

func NewGetPageTemplateQuery(ch *cache.Cache[*template.Template]) *GetPageTemplateQuery {
	return &GetPageTemplateQuery{
		cache: ch,
	}
}

func (q *GetPageTemplateQuery) GetPageTemplate(config *configs.Config, pagename string) (tmpl *template.Template, found bool, err error) {
	key := fmt.Sprintf("%s/%d/%s", config.ID.Hex(), config.UpdatedAt.Unix(), pagename)
	tmpl, ok := q.cache.Get(key)
	if ok {
		return tmpl, true, nil
	}

	page, ok := config.Pages[pagename]
	if !ok {
		return nil, false, nil
	}

	layout, ok := config.Layouts[page.Layout]
	if !ok {
		layout = configs.DefaultLayout
	}

	base, err := configs.PageTemplate.Clone()
	if err != nil {
		return nil, false, fmt.Errorf("failed at cloning page template: %w", err)
	}

	parsed, err := base.Parse(fmt.Sprintf(`{{define "layout_body"}}%s{{end}}{{define "page_body"}}%s{{end}}`, layout.Body, page.Body))
	if err != nil {
		return nil, false, fmt.Errorf("failed at parsing page template body: %w", err)
	}

	q.cache.Set(key, parsed)

	return parsed, true, nil
}
