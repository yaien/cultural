package migrations

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

type Migration struct {
	Name string
	Up   func(ctx context.Context, db *gorm.DB) error
	Down func(ctx context.Context, db *gorm.DB) error
}

var migrations []Migration

func Register(migration Migration) {
	migrations = append(migrations, migration)
}

type Entry struct {
	ID        primitive.ID `gorm:"primaryKey,autoIncrement"`
	Name      string
	AppliedAt time.Time
}

func (Entry) TableName() string {
	return "migrations"
}

func Migrate(ctx context.Context, db *gorm.DB) error {

	if err := db.AutoMigrate(&Entry{}); err != nil {
		return fmt.Errorf("failed migrating migrations table: %w", err)
	}

	slices.SortFunc(migrations, func(a, b Migration) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, migration := range migrations {

		var entry Entry
		err := db.WithContext(ctx).Where("name = ?", migration.Name).Take(&entry).Error
		switch {
		case err == nil:
			continue
		case errors.Is(err, gorm.ErrRecordNotFound):
			err = migration.Up(ctx, db)
			if err != nil {
				return fmt.Errorf("failed applying migration %s: %w", migration.Name, err)
			}

			slog.Info("applied migration", "name", migration.Name)

			entry := &Entry{
				Name:      migration.Name,
				AppliedAt: time.Now(),
			}

			if err := db.WithContext(ctx).Create(entry).Error; err != nil {
				return fmt.Errorf("failed recording migration %s: %w", migration.Name, err)
			}

		default:
			return fmt.Errorf("failed checking migration %s: %w", migration.Name, err)
		}
	}

	return nil

}

func Revert(ctx context.Context, name string, db *gorm.DB) error {
	var entry Entry

	if err := db.Where("name ilike ?", name+"%").Take(&entry).Error; err != nil {
		return fmt.Errorf("failed finding migration %s: %w", name, err)
	}

	var migration *Migration
	for _, m := range migrations {
		if strings.HasSuffix(m.Name, name) {
			migration = &m
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found", name)
	}

	if err := migration.Down(ctx, db); err != nil {
		return fmt.Errorf("failed downgrading migration %s: %w", migration.Name, err)
	}

	if err := db.WithContext(ctx).Delete(&Entry{}, "name = ?", migration.Name).Error; err != nil {
		return fmt.Errorf("failed removing migration record %s: %w", migration.Name, err)
	}

	return nil
}
