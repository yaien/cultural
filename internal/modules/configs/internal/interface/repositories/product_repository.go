package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	db *mongo.Database
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) collection() *mongo.Collection {
	return r.db.Collection("products")
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	_, err := r.collection().InsertOne(ctx, product)
	return err
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}
	_, err := r.collection().UpdateOne(ctx, filter, update)
	return err
}

func (r *ProductRepository) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*models.Product, error) {
	filter := bson.M{"organizationId": organizationID}
	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*models.Product, error) {
	filter := bson.M{"_id": id, "organizationId": organizationID}
	var product models.Product
	err := r.collection().FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, translate(err)
	}
	return &product, nil
}

func (r *ProductRepository) GetBySlugAndOrganizationID(ctx context.Context, slug string, organizationID primitive.ObjectID) (*models.Product, error) {
	filter := bson.M{"slug": slug, "organizationId": organizationID}
	var product models.Product
	err := r.collection().FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, translate(err)
	}
	return &product, nil
}
