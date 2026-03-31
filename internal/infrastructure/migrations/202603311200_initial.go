package migrations

import (
	"context"

	"github.com/yaien/cultural/internal/application/admin"
	"github.com/yaien/cultural/internal/application/auth"
	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/application/store"
	"github.com/yaien/cultural/internal/lib/worker"
	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "202603311200_initial",
		Up: func(ctx context.Context, db *gorm.DB) error {
			return db.AutoMigrate(
				&auth.User{},
				&admin.Group{},
				&admin.Organization{},
				&admin.Invitation{},
				&admin.Role{},
				&label.Config{},
				&label.Draft{},
				&label.Font{},
				&storage.File{},
				&integration.Integration[any]{},
				&store.Product{},
				&store.Presentation{},
				&store.Content{},
				&worker.Job{},
			)
		},
		Down: func(ctx context.Context, db *gorm.DB) error {
			return nil
		},
	})
}
