package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"github.com/yaien/cultural/internal/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.ConfigRepostory = (*ConfigRepository)(nil)

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
		return nil, shared.NotFoundError(err)
	default:
		return nil, err
	}
}
