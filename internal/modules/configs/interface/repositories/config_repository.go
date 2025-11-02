package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.ConfigRepository = (*ConfigRepository)(nil)

type ConfigRepository struct {
	db *mongo.Database
}

func NewConfigRepository(db *mongo.Database) *ConfigRepository {
	return &ConfigRepository{
		db: db,
	}
}

func (r *ConfigRepository) GetByHost(ctx context.Context, host string) (*models.Config, error) {
	var config models.Config
	err := r.db.Collection("configs").FindOne(ctx, bson.M{"host": host}).Decode(&config)
	switch err {
	case nil:
		return &config, nil
	case mongo.ErrNoDocuments:
		return nil, models.NotFoundError(err)
	default:
		return nil, err
	}
}

func (r *ConfigRepository) GetByOrganizationID(ctx context.Context, organizationId primitive.ObjectID) (*models.Config, error) {
	var config models.Config
	err := r.db.Collection("configs").FindOne(ctx, bson.M{"organizationId": organizationId}).Decode(&config)
	switch err {
	case nil:
		return &config, nil
	case mongo.ErrNoDocuments:
		return nil, models.NotFoundError(err)
	default:
		return nil, err
	}
}
