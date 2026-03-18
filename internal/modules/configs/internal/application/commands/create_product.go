package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateProductCommand struct {
	products models.ProductRepository
}

func NewCreateProductCommand(products models.ProductRepository) *CreateProductCommand {
	return &CreateProductCommand{products}
}

type CreateProductRequest struct {
	OrganizationID primitive.ObjectID
	Name           string
}

func (c *CreateProductCommand) CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error) {

	product := &models.Product{
		ID:             primitive.NewObjectID(),
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Slug:           slug.Make(req.Name),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := c.products.GetBySlugAndOrganizationID(ctx, product.Slug, product.OrganizationID)
	switch {
	case err == nil:
		return nil, &models.Error{Code: "product_already_exists", Err: fmt.Errorf("there is already a product with name %q", product.Name)}
	case !models.IsNotFoundError(err):
		return nil, fmt.Errorf("failed at products get by slug and organization id: %w", err)

	}

	if err := c.products.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}
