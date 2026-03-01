package migrations

import (
	"context"

	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrations.Register(migrations.Migration{
		Name: "202602251100_initial_layouts",
		Up: func(ctx context.Context, db *mongo.Database) error {
			cursor, err := db.Collection("configs").Find(ctx, bson.M{})
			if err != nil {
				return err
			}

			var configs []models.Config
			if err := cursor.All(ctx, &configs); err != nil {
				return err
			}

			for _, config := range configs {
				config.Layouts = make(map[string]*models.Layout)
				config.Layouts["default"] = models.DefaultLayout
				for _, page := range config.Pages {
					page.Layout = "default"
				}

				if _, err := db.Collection("configs").ReplaceOne(ctx, bson.M{"_id": config.ID}, config); err != nil {
					return err
				}
			}

			if err = cursor.Close(ctx); err != nil {
				return err
			}

			cursor, err = db.Collection("drafts").Find(ctx, bson.M{})
			if err != nil {
				return err
			}

			var drafts []models.Draft
			if err := cursor.All(ctx, &drafts); err != nil {
				return err
			}

			for _, draft := range drafts {
				draft.Layouts = make(map[string]*models.Layout)
				draft.Layouts["default"] = models.DefaultLayout
				for _, page := range draft.Pages {
					page.Layout = "default"
				}

				if _, err := db.Collection("drafts").ReplaceOne(ctx, bson.M{"_id": draft.ID}, draft); err != nil {
					return err
				}
			}

			if err = cursor.Close(ctx); err != nil {
				return err
			}

			return nil

		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			return nil
		},
	})

}
