package admin

import (
	"context"

	"github.com/yaien/cultural/internal/library/coderror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ RoleRepository = (*MongoRoles)(nil)

type MongoRoles struct {
	collection *mongo.Collection
}

func NewMongoRoles(db *mongo.Database) *MongoRoles {
	return &MongoRoles{db.Collection("roles")}
}

func (r *MongoRoles) CountAdminsByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"organizationId": organizationID, "permissions": "*"})
}

func (r *MongoRoles) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Role, error) {
	var role Role
	err := r.collection.FindOne(ctx, bson.M{"_id": id, "organizationId": organizationID}).Decode(&role)
	switch err {
	case nil:
		return &role, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *MongoRoles) GetByUserIDAndOrganizationID(ctx context.Context, userID, organizationID primitive.ObjectID) (*Role, error) {
	var role Role
	err := r.collection.FindOne(ctx, bson.M{"userId": userID, "organizationId": organizationID}).Decode(&role)
	switch err {
	case nil:
		return &role, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *MongoRoles) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (roles []*Role, err error) {

	cursor, err := r.collection.Find(ctx, bson.M{"organizationId": organizationID})
	if err != nil {
		return nil, err
	}

	defer func() {
		if derr := cursor.Close(ctx); derr != nil {
			err = derr
		}
	}()

	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, err
	}

	return roles, nil

}

func (r *MongoRoles) Create(ctx context.Context, role *Role) error {
	res, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		return err
	}

	role.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (r *MongoRoles) Update(ctx context.Context, role *Role) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": role.ID}, bson.M{"$set": role})
	return err
}

func (r *MongoRoles) Delete(ctx context.Context, role *Role) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": role.ID})
	return err
}
