package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	type File struct {
		ID     primitive.ObjectID `bson:"_id"`
		Preset string             `bson:"preset"`
	}

	type Presentation struct {
		ID       primitive.ObjectID `bson:"_id"`
		Files    []*File            `bson:"files,omitempty"`
		Name     string             `bson:"name"`
		Quantity int                `bson:"quantity"`
		Price    float64            `bson:"price"`
	}

	Register(Migration{
		Name: "202603291620_products_v2",
		Up: func(ctx context.Context, db *mongo.Database) (err error) {
			var products []struct {
				ID            primitive.ObjectID `bson:"_id"`
				Presentations []struct {
					ID       primitive.ObjectID   `bson:"_id"`
					Name     string               `bson:"name"`
					Quantity int                  `bson:"quantity"`
					Price    float64              `bson:"price"`
					FileIDS  []primitive.ObjectID `bson:"fileIds"`
				} `bson:"presentations"`
			}

			cursor, err := db.Collection("products").Find(ctx, bson.M{})
			if err != nil {
				return fmt.Errorf("failed to fetch products: %w", err)
			}

			defer func() {
				if derr := cursor.Close(ctx); derr != nil {
					err = fmt.Errorf("failed to close cursor: %w", derr)
				}
			}()

			if err := cursor.All(ctx, &products); err != nil {
				return fmt.Errorf("failed to decode products: %w", err)
			}

			for _, p := range products {
				var presentations []Presentation
				for _, pres := range p.Presentations {

					files := make([]*File, len(pres.FileIDS))
					for i, id := range pres.FileIDS {
						var file storage.File
						err := db.Collection("files").FindOne(ctx, bson.M{"_id": id}).Decode(&file)
						if err != nil {
							return fmt.Errorf("failed to fetch file with id %s: %w", id.Hex(), err)
						}

						files[i] = &File{ID: file.ID, Preset: file.Preset}
					}

					presentations = append(presentations, Presentation{
						ID:       pres.ID,
						Name:     pres.Name,
						Quantity: pres.Quantity,
						Price:    pres.Price,
						Files:    files,
					})
				}

				_, err := db.Collection("products").UpdateOne(ctx, bson.M{"_id": p.ID}, bson.M{
					"$set": bson.M{
						"presentations": presentations,
						"updatedAt":     time.Now(),
					},
				})

				if err != nil {
					return fmt.Errorf("failed to update product with id %s: %w", p.ID.Hex(), err)
				}

			}

			return nil

		},
	})
}
