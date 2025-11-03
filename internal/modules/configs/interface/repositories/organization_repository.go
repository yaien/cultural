package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.OrganizationRepository = (*OrganizationRepository)(nil)

type OrganizationRepository struct {
	db *mongo.Database
}

func NewOrganizationRepository(db *mongo.Database) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Organization, error) {
	var organization models.Organization
	err := r.db.Collection("organizations").FindOne(ctx, bson.M{"_id": id}).Decode(&organization)
	return &organization, translate(err)
}
