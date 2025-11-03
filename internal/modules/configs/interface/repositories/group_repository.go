package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.GroupRepository = (*GroupRepository)(nil)

type GroupRepository struct {
	db *mongo.Database
}

func NewGroupRepository(db *mongo.Database) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) GetByIDAndOrganizationID(ctx context.Context, id, organizationId primitive.ObjectID) (*models.Group, error) {
	var group models.Group
	err := r.db.Collection("groups").FindOne(ctx, bson.M{"id": id, "organizationId": organizationId}).Decode(&group)
	return &group, translate(err)
}
