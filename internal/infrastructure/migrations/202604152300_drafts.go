package migrations

import (
	"context"

	"github.com/yaien/cultural/internal/application/label"
	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "202604152300_drafts",
		Up: func(ctx context.Context, db *gorm.DB) error {
			configs, err := gorm.G[label.Config](db).Find(ctx)
			if err != nil {
				return err
			}

			for _, config := range configs {
				if len(config.Layouts) == 0 {
					config.Layouts = label.DefaultLayouts
				}

				draft := label.Draft{
					ConfigID: config.ID,
					Layouts:  config.Layouts,
					Fonts:    config.Fonts,
					Pages:    config.Pages,
					Emails:   config.Emails,
					Colors:   config.Colors,
				}

				err := gorm.G[label.Draft](db).Create(ctx, &draft)
				if err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(ctx context.Context, db *gorm.DB) error {
			return nil
		},
	})
}
