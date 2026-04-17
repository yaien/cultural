package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"

	"github.com/yaien/cultural/internal/lib/coderror"
)

type Presentations struct {
	products gorm.Interface[Product]
}

func NewPresentations(db *gorm.DB) *Presentations {
	return &Presentations{products: gorm.G[Product](db)}
}

type CreatePresentationOptions struct {
	ProductID      primitive.ID
	OrganizationID primitive.ID
}

// Create adds a new presentation to the specified product. It retrieves the product by its ID and organization ID,
// creates a new presentation with default values, appends it to the product's presentations, and updates the product in the repository.
// It returns the updated product and the newly created presentation.
func (c *Presentations) Create(ctx context.Context, req *CreatePresentationOptions) (*Product, *Presentation, error) {

	product, err := c.products.
		Where("id = ? and organization_id = ?", req.ProductID, req.OrganizationID).
		First(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	presentation := Presentation{
		ID:   primitive.NewUUID(),
		Name: "Nueva",
	}

	product.Presentations = append(product.Presentations, presentation)

	if _, err = c.products.Updates(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product with new presentation: %w", err)
	}

	return &product, &presentation, nil
}

type UpdatePresentationOptions struct {
	PresentationID primitive.UUID
	ProductID      primitive.ID
	OrganizationID primitive.ID
	Name           string
	Quantity       int
	Price          float64
}

// Update modifies an existing presentation of a product. It validates the input fields, retrieves the product by its ID and organization ID,
func (c *Presentations) Update(ctx context.Context, req *UpdatePresentationOptions) (*Product, *Presentation, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, nil, coderror.New("invalid_presentation_name", fmt.Errorf("presentation name cannot be empty"))
	}

	if req.Quantity < 0 {
		return nil, nil, coderror.New("invalid_presentation_quantity", fmt.Errorf("presentation quantity cannot be less than zero"))
	}

	if req.Price < 0 {
		return nil, nil, coderror.New("invalid_presentation_price", fmt.Errorf("presentation price cannot be less than zero"))
	}

	product, err := c.products.
		Where("id = ? and organization_id = ?", req.ProductID, req.OrganizationID).
		First(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	var updated *Presentation
	for _, presentation := range product.Presentations {
		if presentation.ID == req.PresentationID {
			presentation.Name = req.Name
			presentation.Quantity = req.Quantity
			presentation.Price = req.Price
			updated = &presentation
			break
		}
	}

	if updated == nil {
		return nil, nil, coderror.Newf("presentation_not_found", "presentation with id %s not found", req.PresentationID)
	}

	if _, err = c.products.Updates(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation: %w", err)
	}

	return &product, updated, nil
}

type DeletePresentationOptions struct {
	ID             primitive.UUID
	ProductID      primitive.ID
	OrganizationID primitive.ID
}

// Delete removes a presentation from a product. It retrieves the product by its ID and organization ID, finds the presentation by its ID,
func (c *Presentations) Delete(ctx context.Context, req *DeletePresentationOptions) (*Product, error) {
	product, err := c.products.
		Where("id = ? and organization_id = ?", req.ProductID, req.OrganizationID).
		First(ctx)

	if err != nil {
		return nil, fmt.Errorf("error fetching product: %w", err)
	}

	var deleted *Presentation
	for i, presentation := range product.Presentations {
		if presentation.ID == req.ID {
			product.Presentations = append(product.Presentations[:i], product.Presentations[i+1:]...)
			deleted = &presentation
			break
		}
	}

	if deleted == nil {
		return nil, coderror.Newf("presentation_not_found", "presentation with id %s not found", req.ID)
	}

	if _, err = c.products.Updates(ctx, product); err != nil {
		return nil, fmt.Errorf("error deleting product presentation: %w", err)
	}

	return &product, nil
}
