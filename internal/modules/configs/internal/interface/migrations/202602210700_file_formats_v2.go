package migrations

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrations.Register(migrations.Migration{
		Name: "202602210700_file_formats_v2",
		Up: func(ctx context.Context, db *mongo.Database) error {
			type OldFormat struct {
				ID      primitive.ObjectID `bson:"_id"`
				Quality int                `bson:"quality"`
				Size    int64              `bson:"size"`
				Width   int                `bson:"width"`
				Height  int                `bson:"height"`
			}

			type OldFile struct {
				ID             primitive.ObjectID `bson:"_id"`
				OrganizationID primitive.ObjectID `bson:"organizationId"`
				Name           string             `bson:"name"`
				Formats        map[int]OldFormat  `bson:"formats"`
				ContentType    string             `bson:"contentType"`
				CreatedAt      time.Time          `bson:"createdAt"`
				UpdatedAt      time.Time          `bson:"updatedAt"`
			}

			type NewFormat struct {
				ID          primitive.ObjectID `bson:"_id"`
				Variant     int                `bson:"variant"`
				Size        int64              `bson:"size"`
				Width       int                `bson:"width"`
				Height      int                `bson:"height"`
				ContentType string             `bson:"contentType"`
			}

			type NewFile struct {
				ID             primitive.ObjectID `bson:"_id"`
				OrganizationID primitive.ObjectID `bson:"organizationId"`
				Name           string             `bson:"name"`
				Formats        map[int]NewFormat  `bson:"formats"`
				CreatedAt      time.Time          `bson:"createdAt"`
				UpdatedAt      time.Time          `bson:"updatedAt"`
			}

			var files []OldFile
			cursor, err := db.Collection("files").Find(ctx, bson.M{})
			if err != nil {
				return err
			}
			if err := cursor.All(ctx, &files); err != nil {
				return err
			}

			for _, oldFile := range files {
				newFormats := make(map[int]NewFormat)
				for variant, oldFormat := range oldFile.Formats {
					newFormats[variant] = NewFormat{
						ID:          oldFormat.ID,
						Variant:     oldFormat.Quality,
						Size:        oldFormat.Size,
						Width:       oldFormat.Width,
						Height:      oldFormat.Height,
						ContentType: oldFile.ContentType,
					}
				}

				newFile := NewFile{
					ID:             oldFile.ID,
					OrganizationID: oldFile.OrganizationID,
					Name:           oldFile.Name,
					Formats:        newFormats,
					CreatedAt:      oldFile.CreatedAt,
					UpdatedAt:      oldFile.UpdatedAt,
				}

				_, err := db.Collection("files").ReplaceOne(ctx, bson.M{"_id": oldFile.ID}, newFile)
				if err != nil {
					return err
				}
			}

			return nil

		},
		Down: func(ctx context.Context, db *mongo.Database) error {
			return nil
		},
	})
}
