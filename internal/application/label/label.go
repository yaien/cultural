package label

import (
	"github.com/yaien/cultural/internal/lib/cache"
	"gorm.io/gorm"
)

type Cache = cache.Cache[*Config]

type Label struct {
	Fonts   *Fonts
	Configs *Configs
	Drafts  *Drafts
}

func New(db *gorm.DB, ch *Cache) *Label {
	return &Label{
		Fonts:   NewFonts(db),
		Configs: NewConfigs(db),
		Drafts:  NewDrafts(db, ch),
	}
}
