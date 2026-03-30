package label

import (
	"context"

	"github.com/yaien/cultural/internal/coderror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ ConfigRepository = (*MongoConfigs)(nil)

type MongoConfigs struct {
	collection *mongo.Collection
}

func NewMongoConfigs(db *mongo.Database) *MongoConfigs {
	return &MongoConfigs{db.Collection("configs")}
}

func (r *MongoConfigs) GetByHost(ctx context.Context, host string) (*Config, error) {
	var config Config
	err := r.collection.FindOne(ctx, bson.M{"host": host}).Decode(&config)
	switch err {
	case nil:
		return &config, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *MongoConfigs) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (*Config, error) {
	var config Config
	err := r.collection.FindOne(ctx, bson.M{"organizationId": organizationID}).Decode(&config)
	switch err {
	case nil:
		return &config, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *MongoConfigs) Update(ctx context.Context, config *Config) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": config.ID}, bson.M{"$set": config})
	return err
}
