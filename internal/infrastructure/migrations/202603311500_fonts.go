package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/infrastructure"
	"gorm.io/gorm"

	"google.golang.org/api/option"
	"google.golang.org/api/webfonts/v1"
)

func init() {
	Register(Migration{
		Name: "202603311500_fonts",
		Up: func(ctx context.Context, db *gorm.DB) error {

			config := infrastructure.LoadConfig()
			srv, err := webfonts.NewService(ctx, option.WithAPIKey(config.Google.APIKey))
			if err != nil {
				return fmt.Errorf("failed creating webfonts service: %w", err)
			}

			list, err := srv.Webfonts.List().Capability("WOFF2").Context(ctx).Do()
			if err != nil {
				return fmt.Errorf("failed fetching google fonts: %w", err)
			}

			fonts := make([]*label.Font, len(list.Items))
			for index, item := range list.Items {
				fonts[index] = &label.Font{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Provider:  "google",
					Family:    item.Family,
					Category:  item.Category,
					Subsets:   item.Subsets,
					Variants:  item.Variants,
					Version:   item.Version,
					Files:     item.Files,
				}
			}

			if err := db.WithContext(ctx).Create(fonts); err != nil {
				return fmt.Errorf("failed inserting google fonts: %w", err)
			}

			return nil

		},
		Down: func(ctx context.Context, db *gorm.DB) error {
			return nil
		},
	})
}
