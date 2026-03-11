package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.IntegrationRepository = (*IntegrationRepository)(nil)

type IntegrationRepository struct {
	collection *mongo.Collection
}

func NewIntegrationRepository(db *mongo.Database) *IntegrationRepository {
	return &IntegrationRepository{db.Collection("integrations")}
}

func (i *IntegrationRepository) Create(ctx context.Context, integration *models.Integration) error {
	_, err := i.collection.InsertOne(ctx, integration)
	return err
}

func (i *IntegrationRepository) Update(ctx context.Context, integration *models.Integration) error {
	_, err := i.collection.UpdateOne(ctx, bson.M{"_id": integration.ID}, bson.M{"$set": integration})
	return err
}

func (i *IntegrationRepository) Get(ctx context.Context, options models.GetIntegrationOptions) (*models.Integration, error) {
	var integration models.Integration
	integration.Data = options.Data

	err := i.collection.FindOne(ctx, bson.M{"organizationId": options.OrganizationID, "name": options.Name}).Decode(&integration)
	if err != nil {
		return nil, translate(err)
	}
	return &integration, nil
}
