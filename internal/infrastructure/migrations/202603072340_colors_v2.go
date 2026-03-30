package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	Register(Migration{
		Name: "202603071340_colors_v2",
		Up: func(ctx context.Context, db *mongo.Database) error {
			type Config struct {
				ID     primitive.ObjectID `bson:"_id,omitempty"`
				Colors map[string]string  `bson:"colors"`
			}

			type Color struct {
				ID    primitive.ObjectID `bson:"_id,omitempty"`
				Tag   string             `bson:"tag"`
				Value string             `bson:"value"`
			}

			var configs []Config
			cursor, err := db.Collection("configs").Find(ctx, bson.M{})
			if err != nil {
				return err
			}

			if err := cursor.All(ctx, &configs); err != nil {
				return err
			}

			for _, config := range configs {
				var colors []Color
				for tag, value := range config.Colors {
					color := Color{ID: primitive.NewObjectID(), Tag: tag, Value: value}
					colors = append(colors, color)
				}

				update := bson.M{"$set": bson.M{"colors": colors}}

				if _, err := db.Collection("configs").UpdateOne(ctx, bson.M{"_id": config.ID}, update); err != nil {
					return err
				}

				if _, err := db.Collection("drafts").UpdateOne(ctx, bson.M{"configId": config.ID}, update); err != nil {
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
