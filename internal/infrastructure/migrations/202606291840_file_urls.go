package migrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/yaien/cultural/internal/application/label"
	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "202606291840_file_urls",
		Up: func(ctx context.Context, db *gorm.DB) error {

			var configs []label.Config
			err := db.Find(&configs).Error
			if err != nil {
				return fmt.Errorf("failed getting configs")
			}

			for _, config := range configs {
				for _, page := range config.Pages {
					page.Body = strings.ReplaceAll(page.Body, ".FileURL", "file_url")
					page.Body = strings.ReplaceAll(page.Body, ".ExternalFileURL", "external_file_url")
				}

				for _, layout := range config.Layouts {
					layout.Body = strings.ReplaceAll(layout.Body, ".FileURL", "file_url")
					layout.Body = strings.ReplaceAll(layout.Body, ".ExternalFileURL", "external_file_url")
				}

				err := db.Save(&config).Error
				if err != nil {
					return fmt.Errorf("failed saving config")
				}
			}

			var drafts []label.Draft
			err = db.Find(&drafts).Error
			if err != nil {
				return fmt.Errorf("failed getting drafts: %w", err)
			}

			for _, draft := range drafts {
				for _, page := range draft.Pages {
					page.Body = strings.ReplaceAll(page.Body, ".FileURL", "file_url")
					page.Body = strings.ReplaceAll(page.Body, ".ExternalFileURL", "external_file_url")

				}

				for _, layout := range draft.Layouts {
					layout.Body = strings.ReplaceAll(layout.Body, ".FileURL", "file_url")
					layout.Body = strings.ReplaceAll(layout.Body, ".ExternalFileURL", "external_file_url")
				}

				err := db.Save(&draft).Error
				if err != nil {
					return fmt.Errorf("failed saving draft")
				}
			}

			return nil
		},
	})
}
