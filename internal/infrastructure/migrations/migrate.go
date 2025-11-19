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

func Revert(ctx context.Context, name string, db *mongo.Database) error {
	collection := db.Collection("migrations")
	var entry MigrationEntry

	err := collection.FindOne(ctx, bson.M{"name": bson.M{"$regex": name + "$"}}).Decode(&entry)
	if err != nil {
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

	err = migration.Down(ctx, db)
	if err != nil {
		return fmt.Errorf("failed downgrading migration %s: %w", migration.Name, err)
	}

	slog.Info("downgraded migration", "name", name)

	_, err = collection.DeleteOne(ctx, bson.M{"name": migration.Name})
	if err != nil {
		return fmt.Errorf("failed removing migration record %s: %w", migration.Name, err)
	}

	return nil
}
