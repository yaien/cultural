package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
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

func (r *RoleRepository) GetByUserIDAndOrganizationID(ctx context.Context, userId, organizationId primitive.ObjectID) (*models.Role, error) {
	var role models.Role

	err := r.db.Collection("roles").FindOne(ctx, bson.M{"userId": userId, "organizationId": organizationId}).Decode(&role)

	switch err {
	case nil:
		return &role, nil
	case mongo.ErrNoDocuments:
		return nil, models.NotFoundError(err)
	default:
		return nil, err
	}

}
