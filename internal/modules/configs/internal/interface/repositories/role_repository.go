package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.RoleRepository = (*RoleRepository)(nil)

type RoleRepository struct {
	db *mongo.Database
}

func NewRoleRepository(db *mongo.Database) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) CountAdminsByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (int64, error) {
	count, err := r.db.Collection("roles").CountDocuments(ctx, bson.M{"organizationId": organizationID, "permissions": "*"})
	return count, translate(err)
}

func (r *RoleRepository) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*models.Role, error) {
	var role models.Role
	err := r.db.Collection("roles").FindOne(ctx, bson.M{"_id": id, "organizationId": organizationID}).Decode(&role)
	return &role, translate(err)
}

func (r *RoleRepository) GetByUserIDAndOrganizationID(ctx context.Context, userID, organizationID primitive.ObjectID) (*models.Role, error) {
	var role models.Role
	err := r.db.Collection("roles").FindOne(ctx, bson.M{"userId": userID, "organizationId": organizationID}).Decode(&role)
	return &role, translate(err)
}

func (r *RoleRepository) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*models.Role, error) {
	var roles []*models.Role

	cursor, err := r.db.Collection("roles").Find(ctx, bson.M{"organizationId": organizationID})
	if err != nil {
		return nil, translate(err)
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, translate(err)
	}

	return roles, nil

}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	res, err := r.db.Collection("roles").InsertOne(ctx, role)
	if err != nil {
		return err
	}

	role.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	_, err := r.db.Collection("roles").UpdateOne(ctx, bson.M{"_id": role.ID}, bson.M{"$set": role})
	return err
}

func (r *RoleRepository) Delete(ctx context.Context, role *models.Role) error {
	_, err := r.db.Collection("roles").DeleteOne(ctx, bson.M{"_id": role.ID})
	return err
}
