package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"github.com/yaien/cultural/internal/library/worker"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrations.Register(migrations.Migration{
		Name: "202602221535_file_formats_queue",
		Up: func(ctx context.Context, db *mongo.Database) error {
			if _, err := db.Collection(worker.DefaultJobsCollection).DeleteMany(ctx, bson.M{"name": models.GenerateFormatsTaskName}); err != nil {
				return fmt.Errorf("failed to delete old jobs: %w", err)
			}

			var files []models.File
			cursor, err := db.Collection("files").Find(ctx, bson.M{"format": bson.M{"$exists": false}})
			if err != nil {
				return fmt.Errorf("failed to find files without format: %w", err)
			}
			defer func() {
				if err = cursor.Close(ctx); err != nil {
					slog.Error("Failed to close cursor", "error", err)
				}
			}()

			err = cursor.All(ctx, &files)
			if err != nil {
				return fmt.Errorf("failed to decode files: %w", err)
			}

			queue := worker.NewQueue(worker.NewMongoStore(db, ""), nil)
			for _, file := range files {
				if err := queue.Push(ctx, file.NewGenerateFormatsTask()); err != nil {
					return fmt.Errorf("failed to push job for file %s: %w", file.ID.Hex(), err)
				}
			}

			return nil
		},
	})
}
