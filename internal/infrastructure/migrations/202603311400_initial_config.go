package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"gorm.io/gorm"

	"github.com/yaien/cultural/internal/application/admin"
	"github.com/yaien/cultural/internal/application/label"
)

func init() {
	Register(Migration{
		Name: "202603311400_initial_config",
		Up: func(ctx context.Context, db *gorm.DB) error {
			organization := &admin.Organization{
				Name:      "Cultural",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := db.WithContext(ctx).Create(organization).Error; err != nil {
				return fmt.Errorf("failed to create organization: %w", err)
			}

			cfg := infrastructure.LoadConfig()
			config := &label.Config{
				Host:           cfg.Init.Host,
				Title:          cfg.Init.Title,
				Url:            cfg.Init.Url,
				Email:          cfg.Init.Email,
				OrganizationID: organization.ID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				Colors:         label.DefaultColors,
				Pages:          label.DefaultPages,
				Emails:         label.DefaultEmails,
				Layouts:        label.DefaultLayouts,
			}

			if err := db.WithContext(ctx).Create(config).Error; err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}

			return nil

		},
		Down: func(ctx context.Context, db *gorm.DB) error {
			return nil
		},
	})
}
