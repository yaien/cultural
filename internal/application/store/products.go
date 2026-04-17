package store

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"

	"github.com/gosimple/slug"
	"github.com/yaien/cultural/internal/lib/coderror"
)

type Products struct {
	products gorm.Interface[Product]
}

func NewProducts(db *gorm.DB) *Products {
	return &Products{gorm.G[Product](db)}
}

type CreateProductOptions struct {
	OrganizationID primitive.ID
	Name           string
}

// Create creates a new product with the given options.
// It checks if a product with the same slug already exists for the organization and returns an error if it does.
// Otherwise, it creates the product in the repository and returns it.
func (c *Products) Create(ctx context.Context, req *CreateProductOptions) (*Product, error) {

	product := &Product{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Slug:           slug.Make(req.Name),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if req.Name == "" {
		return nil, coderror.New("invalid_product_name", fmt.Errorf("product name cannot be empty"))
	}

	count, err := c.products.Where("slug = ? and organization_id = ?", product.Slug, product.OrganizationID).Count(ctx, "id")
	if err != nil {
		return nil, primitive.Error(fmt.Errorf("failed at products get by slug and organization id: %w", err))
	}

	if count > 0 {
		return nil, coderror.Newf("product_already_exists", "there is already a product with name %q", product.Name)
	}

	if err := c.products.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetByOrganizationID retrieves all products for a given organization ID.
func (c *Products) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) ([]Product, error) {
	return c.products.
		Where("organization_id = ?", organizationID).
		Find(ctx)
}

// GetByIDAndOrganizationID retrieves a product by its ID and organization ID.
func (c *Products) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (Product, error) {
	return c.products.
		Where("id = ? and organization_id = ?", id, organizationID).
		Take(ctx)
}
