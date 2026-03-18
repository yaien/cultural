package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetProductByIDQuery struct {
	products models.ProductRepository
}

func NewGetProductByIDQuery(products models.ProductRepository) *GetProductByIDQuery {
	return &GetProductByIDQuery{products: products}
}

func (q *GetProductByIDQuery) GetProductByID(ctx context.Context, productID, organizationID primitive.ObjectID) (*models.Product, error) {
	return q.products.GetByIDAndOrganizationID(ctx, productID, organizationID)
}
