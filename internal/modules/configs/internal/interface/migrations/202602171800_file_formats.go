package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"github.com/yaien/cultural/internal/library/worker"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {

	migrations.Register(migrations.Migration{
		Name: "202602171800_file_formats",
		Up: func(ctx context.Context, db *mongo.Database) error {
			type OldFile struct {
				ID             primitive.ObjectID `bson:"_id"`
				OrganizationID primitive.ObjectID `bson:"organizationId"`
				Name           string             `bson:"name"`
				Size           int64              `bson:"size"`
				Formats        []models.Format    `bson:"formats"`
				MimeType       string             `bson:"mimeType"`
				CreatedAt      time.Time          `bson:"createdAt"`
				UpdatedAt      time.Time          `bson:"updatedAt"`
			}

			type NewFile struct {
				ID             primitive.ObjectID    `bson:"_id"`
				OrganizationID primitive.ObjectID    `bson:"organizationId"`
				Name           string                `bson:"name"`
				ContentType    string                `bson:"contentType"`
				Formats        map[int]models.Format `bson:"formats"`
				CreatedAt      time.Time             `bson:"createdAt"`
				UpdatedAt      time.Time             `bson:"updatedAt"`
			}

			cfg := infrastructure.LoadConfig()
			collection := db.Collection("files")
			root := cfg.Storage.Local.Path
			queue := worker.NewQueue(worker.NewMongoStore(db, ""), nil)

			var oldies []OldFile
			cursor, err := collection.Find(ctx, bson.M{})
			if err != nil {
				return fmt.Errorf("failed finding oldies: %w", err)
			}

			if err := cursor.All(ctx, &oldies); err != nil {
				return fmt.Errorf("failed decoding oldies: %w", err)
			}

			for _, oldie := range oldies {

				// Skip files that already have formats (e.g., from previous runs or manual updates)
				if len(oldie.Formats) > 0 {
					continue
				}

				width, height, quality, err := models.GetFileDimension(root, oldie.ID.Hex(), oldie.MimeType)
				if err != nil {
					return fmt.Errorf("failed getting dimensions for file %s: %w", oldie.ID.Hex(), err)
				}

				format := models.Format{
					ID:      oldie.ID,
					Size:    oldie.Size,
					Width:   width,
					Height:  height,
					Variant: quality,
				}

				newbie := &NewFile{
					ID:             oldie.ID,
					OrganizationID: oldie.OrganizationID,
					Name:           oldie.Name,
					ContentType:    oldie.MimeType,
					Formats:        map[int]models.Format{format.Variant: format},
					CreatedAt:      oldie.CreatedAt,
					UpdatedAt:      time.Now(),
				}

				_, err = collection.ReplaceOne(ctx, bson.M{"_id": oldie.ID}, newbie)
				if err != nil {
					return fmt.Errorf("failed updating file: %w", err)
				}

				file := &models.File{ID: oldie.ID}
				if err := queue.Push(ctx, file.NewGenerateFormatsTask()); err != nil {
					return fmt.Errorf("failed pushing optimize task: %w", err)
				}

			}

			return nil

		},
	})
}
