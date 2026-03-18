package migrations

import (
	"context"

	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrations.Register(migrations.Migration{
		Name: "202603171800_missing_indexes",
		Up: func(ctx context.Context, db *mongo.Database) error {

			collection := db.Collection("roles")
			_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "organizationId", Value: 1}},
			})

			if err != nil {
				return err
			}

			draftsCollection := db.Collection("drafts")
			_, err = draftsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{{Key: "configId", Value: 1}},
			})

			if err != nil {
				return err
			}

			filesCollection := db.Collection("files")
			_, err = filesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{{Key: "organizationId", Value: 1}, {Key: "name", Value: 1}},
			})

			if err != nil {
				return err
			}

			usersCollection := db.Collection("users")
			_, err = usersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{{Key: "email", Value: 1}},
			})

			if err != nil {
				return err
			}

			invitationsCollection := db.Collection("invitations")
			_, err = invitationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{{Key: "userEmail", Value: 1}, {Key: "organizationId", Value: 1}},
			})

			if err != nil {
				return err
			}

			integrationsCollection := db.Collection("integrations")
			_, err = integrationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys: bson.D{{Key: "organizationId", Value: 1}},
			})

			return err

		},
	})
}
