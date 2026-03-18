package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetProductsQuery struct {
	products models.ProductRepository
}

func NewGetProductsQuery(products models.ProductRepository) *GetProductsQuery {
	return &GetProductsQuery{products: products}
}

func (q *GetProductsQuery) GetProducts(ctx context.Context, organizationID primitive.ObjectID) ([]*models.Product, error) {
	return q.products.GetByOrganizationID(ctx, organizationID)
}
