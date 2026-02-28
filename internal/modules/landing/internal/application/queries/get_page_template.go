package queries

import (
	"fmt"

	"github.com/yaien/cultural/internal/library/cache"
	"github.com/yaien/cultural/internal/modules/configs"
)

type GetPageTemplateQuery struct {
	cache *cache.Cache[string]
}

func NewGetPageTemplateQuery(ch *cache.Cache[string]) *GetPageTemplateQuery {
	return &GetPageTemplateQuery{
		cache: ch,
	}
}

func (q *GetPageTemplateQuery) GetPageHTML(config *configs.Config, pagename string) (html string, found bool, err error) {
	if pagename == "index" {
		return "", false, nil
	}

	if pagename == "" {
		pagename = "index"
	}

	key := fmt.Sprintf("%s/%d/%s", config.ID.Hex(), config.UpdatedAt.Unix(), pagename)
	html, ok := q.cache.Get(key)
	if ok {
		return html, true, nil
	}

	page, ok := config.Pages[pagename]
	if !ok {
		return "", false, nil
	}

	layout, ok := config.Layouts[page.Layout]
	if !ok {
		layout = configs.DefaultLayout
	}

	html, err = configs.RenderPage(&configs.PageData{
		Page:     page,
		Layout:   layout,
		AppTitle: config.Title,
		Fonts:    config.Fonts,
		Colors:   config.Colors,
		FilePath: "/assets/dynamic/files/",
	})

	if err != nil {
		return "", false, fmt.Errorf("failed at rendering page: %w", err)
	}

	q.cache.Set(key, html)

	return html, true, nil
}
