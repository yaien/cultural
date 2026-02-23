package commands

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteFileCommand struct {
	files   models.FileRepository
	storage storage.Storage
}

func NewDeleteFileCommand(files models.FileRepository, st storage.Storage) *DeleteFileCommand {
	return &DeleteFileCommand{
		files:   files,
		storage: st,
	}
}

func (c *DeleteFileCommand) DeleteFile(ctx context.Context, organizationID primitive.ObjectID, name string) error {
	file, err := c.files.GetByOrganizationIDAndName(ctx, organizationID, name)
	if err != nil {
		return fmt.Errorf("failed to get file from repository: %w", err)
	}

	err = c.files.DeleteByOrganizationIDAndName(ctx, organizationID, name)
	if err != nil {
		return fmt.Errorf("failed to delete file from repository: %w", err)
	}

	for _, format := range file.Formats {
		err = c.storage.Remove(format.ID.Hex())
		if err != nil {
			return fmt.Errorf("failed to delete file from storage: %w", err)
		}
	}

	return nil
}
