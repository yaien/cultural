package store

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/lib/coderror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ Repository = (*Mongo)(nil)

type Mongo struct {
	collection *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{db.Collection("products")}
}

func (r *Mongo) Create(ctx context.Context, product *Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *Mongo) Update(ctx context.Context, product *Product) error {
	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *Mongo) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (products []*Product, err error) {
	filter := bson.M{"organizationId": organizationID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer func() {
		if derr := cursor.Close(ctx); derr != nil {
			err = fmt.Errorf("failed to close cursor: %w", derr)
		}
	}()

	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *Mongo) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Product, error) {
	filter := bson.M{"_id": id, "organizationId": organizationID}
	var product Product
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	switch err {
	case nil:
		return &product, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *Mongo) GetBySlugAndOrganizationID(ctx context.Context, slug string, organizationID primitive.ObjectID) (*Product, error) {
	filter := bson.M{"slug": slug, "organizationId": organizationID}
	var product Product
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	switch err {
	case nil:
		return &product, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}
