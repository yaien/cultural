package commands

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteProductPresentationCommand struct {
	products models.ProductRepository
}

func NewDeleteProductPresentationCommand(products models.ProductRepository) *DeleteProductPresentationCommand {
	return &DeleteProductPresentationCommand{products: products}
}

type DeleteProductPresentationRequest struct {
	ID             primitive.ObjectID
	ProductID      primitive.ObjectID
	OrganizationID primitive.ObjectID
}

func (c *DeleteProductPresentationCommand) DeleteProductPresentation(ctx context.Context, req DeleteProductPresentationRequest) error {
	product, err := c.products.GetByIDAndOrganizationID(ctx, req.ProductID, req.OrganizationID)
	if err != nil {
		return fmt.Errorf("error fetching product: %w", err)
	}

	var deleted *models.Presentation
	for i, presentation := range product.Presentations {
		if presentation.ID == req.ID {
			product.Presentations = append(product.Presentations[:i], product.Presentations[i+1:]...)
			deleted = presentation
			break
		}
	}

	if deleted == nil {
		return &models.Error{Code: "presentation_not_found", Err: fmt.Errorf("presentation with id %s not found", req.ID.Hex())}
	}

	if err = c.products.Update(ctx, product); err != nil {
		return fmt.Errorf("error deleting product presentation: %w", err)
	}

	return nil
}
