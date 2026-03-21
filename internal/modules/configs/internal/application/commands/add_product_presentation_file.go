package commands

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddProductPresentationFileCommand struct {
	upload   *UploadFileCommand
	products models.ProductRepository
}

func NewAddProductPresentationFileCommand(products models.ProductRepository, upload *UploadFileCommand) *AddProductPresentationFileCommand {
	return &AddProductPresentationFileCommand{products: products, upload: upload}
}

type AddProductPresentationFileRequest struct {
	PresentationID primitive.ObjectID
	ProductID      primitive.ObjectID
	OrganizationID primitive.ObjectID
	Name           string
	Size           int64
	ContentType    string
	Data           io.Reader
}

func (c *AddProductPresentationFileCommand) AddProductPresentationFile(ctx context.Context, req AddProductPresentationFileRequest) (*models.Product, *models.Presentation, error) {
	product, err := c.products.GetByIDAndOrganizationID(ctx, req.ProductID, req.OrganizationID)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	var presentation *models.Presentation
	for _, p := range product.Presentations {
		if p.ID == req.PresentationID {
			presentation = p
		}
	}

	if presentation == nil {
		return nil, nil, &models.Error{Code: "presentation_not_found", Err: fmt.Errorf("presentation with id %s not found", req.PresentationID.Hex())}
	}

	if len(presentation.FileIDS) >= 5 {
		return nil, nil, &models.Error{Code: "presentation_file_limit_exceeded", Err: fmt.Errorf("presentation file limit exceeded")}
	}

	file, err := c.upload.UploadFile(ctx, &UploadFileRequest{
		Name:           req.Name,
		Size:           req.Size,
		ContentType:    req.ContentType,
		Data:           req.Data,
		OrganizationID: req.OrganizationID,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("error uploading file: %w", err)
	}

	presentation.FileIDS = append(presentation.FileIDS, file.ID)
	product.UpdatedAt = time.Now()

	if err = c.products.Update(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation files: %w", err)
	}

	return product, presentation, nil
}
