package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/infrastructure"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/option"
	"google.golang.org/api/webfonts/v1"
)

func init() {
	Register(Migration{
		Name: "202511181900_fonts",
		Up: func(ctx context.Context, db *mongo.Database) error {
			collection := db.Collection("fonts")
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{
					{Key: "family", Value: "text"},
					{Key: "provider", Value: 1},
				},
			})

			if err != nil {
				return fmt.Errorf("failed creating indexes: %w", err)
			}

			config := infrastructure.LoadConfig()
			srv, err := webfonts.NewService(ctx, option.WithAPIKey(config.Google.APIKey))
			if err != nil {
				return fmt.Errorf("failed creating webfonts service: %w", err)
			}

			list, err := srv.Webfonts.List().Capability("WOFF2").Context(ctx).Do()
			if err != nil {
				return fmt.Errorf("failed fetching google fonts: %w", err)
			}

			fonts := make([]any, len(list.Items))
			for index, item := range list.Items {
				fonts[index] = &label.Font{
					ID:        primitive.NewObjectID(),
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

			_, err = collection.InsertMany(ctx, fonts)
			if err != nil {
				return fmt.Errorf("failed inserting google fonts: %w", err)
			}

			return nil

		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			return db.Collection("fonts").Drop(ctx)
		},
	})
}
