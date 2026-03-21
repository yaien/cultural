package commands

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateProductPresentationCommand struct {
	products models.ProductRepository
}

func NewCreateProductPresentationCommand(products models.ProductRepository) *CreateProductPresentationCommand {
	return &CreateProductPresentationCommand{products: products}
}

func (C *CreateProductPresentationCommand) CreateProductPresentation(ctx context.Context, productID primitive.ObjectID, organizationID primitive.ObjectID) (*models.Product, *models.Presentation, error) {
	product, err := C.products.GetByIDAndOrganizationID(ctx, productID, organizationID)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	presentation := &models.Presentation{
		ID:       primitive.NewObjectID(),
		FileIDS:  []primitive.ObjectID{},
		Name:     "Nueva",
		Quantity: 0,
		Price:    0.0,
	}

	product.Presentations = append(product.Presentations, presentation)

	if err = C.products.Update(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product with new presentation: %w", err)
	}

	return product, presentation, nil
}
