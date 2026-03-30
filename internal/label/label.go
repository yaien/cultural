package label

import "github.com/yaien/cultural/internal/cache"

type Cache = cache.Cache[*Config]

type Label struct {
	Fonts   *Fonts
	Configs *Configs
	Drafts  *Drafts
}

func New(fonts FontRepository, configs ConfigRepository, drafts DraftRepository, ch *Cache) *Label {
	return &Label{
		Fonts:   NewFonts(fonts),
		Configs: NewConfigs(configs),
		Drafts:  NewDrafts(drafts, configs, fonts, ch),
	}
}
