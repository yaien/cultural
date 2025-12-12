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
	return &config, translate(err)
}

func (r *ConfigRepository) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (*models.Config, error) {
	var config models.Config
	err := r.db.Collection("configs").FindOne(ctx, bson.M{"organizationId": organizationID}).Decode(&config)
	return &config, translate(err)
}

func (r *ConfigRepository) Update(ctx context.Context, config *models.Config) error {
	_, err := r.db.Collection("configs").UpdateOne(ctx, bson.M{"_id": config.ID}, bson.M{"$set": config})
	return err
}
