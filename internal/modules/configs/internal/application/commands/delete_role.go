package commands

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
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

type DeleteRoleRequest struct {
	SessionRole    *models.Role
	TargetRoleID   primitive.ObjectID
	OrganizationID primitive.ObjectID
}

func (c *DeleteRoleCommand) DeleteRole(ctx context.Context, req *DeleteRoleRequest) error {
	if req.SessionRole.ID == req.TargetRoleID {
		return &models.Error{Code: "invalid_deletion", Err: fmt.Errorf("you cannot delete your own role")}
	}

	if !req.SessionRole.Permissions.Has("delete-role") {
		return &models.Error{Code: "permission_denied", Err: fmt.Errorf("you do not have permission to delete roles")}
	}

	count, err := c.roles.CountAdminsByOrganizationID(ctx, req.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to count roles: %w", err)
	}

	if count == 1 {
		return &models.Error{Code: "last_admin_deletion", Err: fmt.Errorf("cannot delete the last admin role in the organization")}
	}

	role, err := c.roles.GetByIDAndOrganizationID(ctx, req.TargetRoleID, req.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	err = c.roles.Delete(ctx, role)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}
