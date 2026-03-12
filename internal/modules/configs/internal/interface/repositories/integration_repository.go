package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.IntegrationRepository[any] = (*IntegrationRepository[any])(nil)

type IntegrationRepository[T any] struct {
	collection *mongo.Collection
}

func NewIntegrationRepository[T any](db *mongo.Database) *IntegrationRepository[T] {
	return &IntegrationRepository[T]{db.Collection("integrations")}
}

func (i *IntegrationRepository[T]) Create(ctx context.Context, integration *models.Integration[T]) error {
	_, err := i.collection.InsertOne(ctx, integration)
	return err
}

func (i *IntegrationRepository[T]) Update(ctx context.Context, integration *models.Integration[T]) error {
	_, err := i.collection.UpdateOne(ctx, bson.M{"_id": integration.ID}, bson.M{"$set": integration})
	return err
}

func (i *IntegrationRepository[T]) Get(ctx context.Context, options models.GetIntegrationOptions) (*models.Integration[T], error) {
	var integration models.Integration[T]
	err := i.collection.FindOne(ctx, bson.M{"organizationId": options.OrganizationID, "name": options.Name}).Decode(&integration)
	if err != nil {
		return nil, translate(err)
	}
	return &integration, nil
}
