package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateProductPresentationCommand struct {
	products models.ProductRepository
}

func NewUpdateProductPresentationCommand(products models.ProductRepository) *UpdateProductPresentationCommand {
	return &UpdateProductPresentationCommand{products: products}
}

type UpdateProductPresentationRequest struct {
	ID             primitive.ObjectID
	Name           string
	Quantity       int
	Price          float64
	ProductID      primitive.ObjectID
	OrganizationID primitive.ObjectID
}

func (c *UpdateProductPresentationCommand) UpdateProductPresentation(ctx context.Context, req UpdateProductPresentationRequest) (*models.Product, *models.Presentation, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, nil, &models.Error{Code: "invalid_presentation_name", Err: fmt.Errorf("presentation name cannot be empty")}
	}

	if req.Quantity < 0 {
		return nil, nil, &models.Error{Code: "invalid_presentation_quantity", Err: fmt.Errorf("presentation quantity cannot be less than zero")}
	}

	if req.Price < 0 {
		return nil, nil, &models.Error{Code: "invalid_presentation_price", Err: fmt.Errorf("presentation price cannot be less than zero")}
	}

	product, err := c.products.GetByIDAndOrganizationID(ctx, req.ProductID, req.OrganizationID)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	var updated *models.Presentation
	for _, presentation := range product.Presentations {
		if presentation.ID == req.ID {
			presentation.Name = req.Name
			presentation.Quantity = req.Quantity
			presentation.Price = req.Price
			updated = presentation
			break
		}
	}

	if updated == nil {
		return nil, nil, &models.Error{Code: "presentation_not_found", Err: fmt.Errorf("presentation with id %s not found", req.ID.Hex())}
	}

	if err = c.products.Update(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation: %w", err)
	}

	return product, updated, nil
}
