package integration

import (
	"context"

	"github.com/yaien/cultural/internal/coderror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ Repository[any] = (*Mongo[any])(nil)

type Mongo[T any] struct {
	collection *mongo.Collection
}

func NewMongo[T any](db *mongo.Database) *Mongo[T] {
	return &Mongo[T]{db.Collection("integrations")}
}

func (i *Mongo[T]) Create(ctx context.Context, integration *Integration[T]) error {
	_, err := i.collection.InsertOne(ctx, integration)
	return err
}

func (i *Mongo[T]) Update(ctx context.Context, integration *Integration[T]) error {
	_, err := i.collection.UpdateOne(ctx, bson.M{"_id": integration.ID}, bson.M{"$set": integration})
	return err
}

func (i *Mongo[T]) GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) (*Integration[T], error) {
	var integration Integration[T]
	err := i.collection.FindOne(ctx, bson.M{"organizationId": organizationID, "name": name}).Decode(&integration)
	switch err {
	case nil:
		return &integration, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (i *Mongo[T]) GetByName(ctx context.Context, name string) (integrations []*Integration[T], err error) {
	cursor, err := i.collection.Find(ctx, bson.M{"name": name})
	if err != nil {
		return nil, err
	}

	defer func() {
		if derr := cursor.Close(ctx); derr != nil {
			err = derr
		}
	}()

	if err := cursor.All(ctx, &integrations); err != nil {
		return nil, err
	}

	return integrations, nil

}
