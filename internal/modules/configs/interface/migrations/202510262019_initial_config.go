package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrations.Register(migrations.Migration{
		Name: "202510262019_initial_config",
		Up: func(ctx context.Context, db *mongo.Database) error {
			organizations := db.Collection("organizations")
			res, err := organizations.InsertOne(ctx, bson.M{
				"name":      "Cultural App",
				"createdAt": time.Now(),
				"updatedAt": time.Now(),
			})

			if err != nil {
				return err
			}

			configs := db.Collection("configs")
			_, err = configs.Indexes().CreateMany(ctx, []mongo.IndexModel{
				{Keys: bson.D{{Key: "host", Value: 1}}, Options: nil},
				{Keys: bson.D{{Key: "organizationId", Value: 1}}, Options: nil},
			})

			if err != nil {
				return fmt.Errorf("failed creating indexes: %w", err)
			}

			_, err = configs.InsertOne(ctx, bson.M{
				"host":           viper.GetString("INIT_CONFIG_HOST"),
				"title":          viper.GetString("INIT_CONFIG_TITLE"),
				"url":            viper.GetString("INIT_CONFIG_URL"),
				"theme":          viper.GetString("INIT_CONFIG_THEME"),
				"organizationId": res.InsertedID,
				"createdAt":      time.Now(),
				"updatedAt":      time.Now(),
			})

			return err

		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			return nil
		},
	})
}
