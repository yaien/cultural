package admin

import (
	"context"

	"github.com/yaien/cultural/internal/lib/coderror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ GroupRepository = (*MongoGroups)(nil)

type MongoGroups struct {
	collection *mongo.Collection
}

func NewMongoGroups(db *mongo.Database) *MongoGroups {
	return &MongoGroups{db.Collection("groups")}
}

func (r *MongoGroups) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Group, error) {
	var group Group
	err := r.collection.FindOne(ctx, primitive.M{"_id": id, "organizationId": organizationID}).Decode(&group)
	switch err {
	case nil:
		return &group, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}
