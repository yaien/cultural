package migrations

import (
	"context"
	"strings"
	"time"

	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrations.Register(migrations.Migration{
		Name: "202602211400_file_formats_v3",
		Up: func(ctx context.Context, db *mongo.Database) error {
			type Format struct {
				ID          primitive.ObjectID `bson:"_id"`
				Variant     int                `bson:"variant"`
				Size        int64              `bson:"size"`
				Width       int                `bson:"width"`
				Height      int                `bson:"height"`
				ContentType string             `bson:"contentType"`
			}

			type OldFile struct {
				ID             primitive.ObjectID `bson:"_id"`
				OrganizationID primitive.ObjectID `bson:"organizationId"`
				Name           string             `bson:"name"`
				Formats        map[int]Format     `bson:"formats"`
				CreatedAt      time.Time          `bson:"createdAt"`
				UpdatedAt      time.Time          `bson:"updatedAt"`
			}

			type NewFile struct {
				ID             primitive.ObjectID `bson:"_id"`
				OrganizationID primitive.ObjectID `bson:"organizationId"`
				Name           string             `bson:"name"`
				Preset         string             `bson:"preset,omitempty"`
				CreatedAt      time.Time          `bson:"createdAt"`
				UpdatedAt      time.Time          `bson:"updatedAt"`
				Formats        []Format           `bson:"formats"`
				Optimized      bool               `bson:"optimized"`
				OptimizedAt    *time.Time         `bson:"optimizedAt,omitempty"`
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
				var newFormats []Format
				for variant, oldFormat := range oldFile.Formats {
					newFormats = append(newFormats, Format{
						Variant:     variant,
						ID:          oldFormat.ID,
						Size:        oldFormat.Size,
						Width:       oldFormat.Width,
						Height:      oldFormat.Height,
						ContentType: oldFormat.ContentType,
					})
				}

				preset := ""
				if len(newFormats) > 0 {
					preset = strings.Split(newFormats[0].ContentType, "/")[0]
				}

				newFile := NewFile{
					ID:             oldFile.ID,
					OrganizationID: oldFile.OrganizationID,
					Name:           oldFile.Name,
					Preset:         preset,
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
