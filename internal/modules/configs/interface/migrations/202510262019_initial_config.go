package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrations.Register(migrations.Migration{
		Name: "202510262019_initial_config",
		Up: func(ctx context.Context, db *mongo.Database) error {
			organizations := db.Collection("organizations")
			res, err := organizations.InsertOne(ctx, models.Organization{
				Name:      "Cultural",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
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

			cfg := infrastructure.LoadConfig()

			_, err = configs.InsertOne(ctx,
				models.Config{
					Host:           cfg.Init.Host,
					Title:          cfg.Init.Title,
					Url:            cfg.Init.Url,
					Email:          cfg.Init.Email,
					OrganizationID: res.InsertedID.(primitive.ObjectID),
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
					Colors:         models.DefaultColors,
					Pages:          models.DefaultPages,
					Emails:         models.DefaultEmails,
				})

			return err

		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			err := db.Collection("configs").Drop(ctx)
			if err != nil {
				return err
			}

			return db.Collection("organizations").Drop(ctx)
		},
	})
}
