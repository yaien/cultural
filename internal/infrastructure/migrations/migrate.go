package migrations

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	Name string
	Up   func(ctx context.Context, db *mongo.Database) error
	Down func(ctx context.Context, db *mongo.Database) error
}

var migrations []Migration

func Register(migration Migration) {
	migrations = append(migrations, migration)
}

type MigrationEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	AppliedAt time.Time `bson:"appliedAt" json:"appliedAt"`
}

func Migrate(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("migrations")

	slog.Info("starting migrations", "count", len(migrations))

	slices.SortFunc(migrations, func(a, b Migration) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, migration := range migrations {
		var entry MigrationEntry

		err := collection.FindOne(ctx, bson.M{"name": migration.Name}).Decode(&entry)
		switch {
		case err == nil:
			continue
		case errors.Is(err, mongo.ErrNoDocuments):
			err = migration.Up(ctx, db)
			if err != nil {
				return fmt.Errorf("failed applying migration %s: %w", migration.Name, err)
			}

			slog.Info("applied migration", "name", migration.Name)

			_, err := collection.InsertOne(ctx, MigrationEntry{
				Name:      migration.Name,
				AppliedAt: time.Now(),
			})

			if err != nil {
				return fmt.Errorf("failed recording migration %s: %w", migration.Name, err)
			}

		default:
			return fmt.Errorf("failed checking migration %s: %w", migration.Name, err)
		}
	}

	return nil

}
