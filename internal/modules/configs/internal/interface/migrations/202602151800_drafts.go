package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrations.Register(migrations.Migration{
		Name: "20260201151800_drafts",
		Up: func(ctx context.Context, db *mongo.Database) error {
			var configs []*models.Config
			cursor, err := db.Collection("configs").Find(ctx, bson.M{})
			if err != nil {
				return fmt.Errorf("failed to find configs: %w", err)
			}

			if err := cursor.All(ctx, &configs); err != nil {
				return fmt.Errorf("failed to decode configs: %w", err)
			}

			for _, config := range configs {
				draft := &models.Draft{
					ID:        primitive.NewObjectID(),
					ConfigID:  config.ID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Layouts:   config.Layouts,
					Fonts:     config.Fonts,
					Pages:     config.Pages,
					Emails:    config.Emails,
					Colors:    config.Colors,
				}

				if _, err := db.Collection("drafts").InsertOne(ctx, draft); err != nil {
					return fmt.Errorf("failed to insert draft for config %s: %w", config.ID.Hex(), err)
				}
			}

			return nil
		},
	})
}
