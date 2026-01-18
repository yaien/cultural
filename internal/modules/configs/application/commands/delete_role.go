package commands

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteRoleCommand struct {
	roles models.RoleRepository
}

func NewDeleteRoleCommand(roles models.RoleRepository) *DeleteRoleCommand {
	return &DeleteRoleCommand{
		roles: roles,
	}
}

func (c *DeleteRoleCommand) DeleteRole(ctx context.Context, userID, organizationID primitive.ObjectID) error {
	role, err := c.roles.GetByUserIDAndOrganizationID(ctx, userID, organizationID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	err = c.roles.Delete(ctx, role)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}
