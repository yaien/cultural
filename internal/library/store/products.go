package store

import (
	"context"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"github.com/yaien/cultural/internal/library/coderror"
	"github.com/yaien/cultural/internal/library/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Products struct {
	repository Repository
	storage    *storage.Storage
}

func NewProducts(repository Repository, storage *storage.Storage) *Products {
	return &Products{repository: repository, storage: storage}
}

type CreateProductOptions struct {
	OrganizationID primitive.ObjectID
	Name           string
}

// Create creates a new product with the given options.
// It checks if a product with the same slug already exists for the organization and returns an error if it does.
// Otherwise, it creates the product in the repository and returns it.
func (c *Products) Create(ctx context.Context, req *CreateProductOptions) (*Product, error) {

	product := &Product{
		ID:             primitive.NewObjectID(),
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Slug:           slug.Make(req.Name),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if req.Name == "" {
		return nil, coderror.New("invalid_product_name", fmt.Errorf("product name cannot be empty"))
	}

	_, err := c.repository.GetBySlugAndOrganizationID(ctx, product.Slug, product.OrganizationID)
	switch {
	case err == nil:
		return nil, coderror.Newf("product_already_exists", "there is already a product with name %q", product.Name)
	case !coderror.Is(err, coderror.NotFound):
		return nil, fmt.Errorf("failed at products get by slug and organization id: %w", err)
	}

	if err := c.repository.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetByOrganizationID retrieves all products for a given organization ID.
func (c *Products) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*Product, error) {
	return c.repository.GetByOrganizationID(ctx, organizationID)
}

// GetByIDAndOrganizationID retrieves a product by its ID and organization ID.
func (c *Products) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Product, error) {
	return c.repository.GetByIDAndOrganizationID(ctx, id, organizationID)
}
