package commands

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RenameFileCommand struct {
	files models.FileRepository
}

func NewRenameFileCommand(files models.FileRepository) *RenameFileCommand {
	return &RenameFileCommand{files: files}
}

func (c *RenameFileCommand) RenameFile(ctx context.Context, organizationId primitive.ObjectID, oldName, newName string) error {
	return c.files.Rename(ctx, organizationId, oldName, newName)
}
