package admin

import (
	"context"

	"github.com/yaien/cultural/internal/library/coderror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ OrganizationRepository = (*MongoOrganizations)(nil)

type MongoOrganizations struct {
	collection *mongo.Collection
}

func NewMongoOrganizations(db *mongo.Database) *MongoOrganizations {
	return &MongoOrganizations{db.Collection("organizations")}
}

func (r *MongoOrganizations) GetByID(ctx context.Context, id primitive.ObjectID) (*Organization, error) {
	var organization Organization
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&organization)
	switch err {
	case nil:
		return &organization, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}
