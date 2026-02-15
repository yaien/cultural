package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateRoleCommand struct {
	roles  models.RoleRepository
	groups models.GroupRepository
}

func NewUpdateRoleCommand(roles models.RoleRepository, groups models.GroupRepository) *UpdateRoleCommand {
	return &UpdateRoleCommand{
		roles:  roles,
		groups: groups,
	}
}

type UpdateRoleRequest struct {
	ID             primitive.ObjectID
	OrganizationID primitive.ObjectID
	GroupID        *primitive.ObjectID
	Permissions    models.Permissions
	Name           string
}

func (c *UpdateRoleCommand) UpdateRole(ctx context.Context, r *UpdateRoleRequest) error {
	role, err := c.roles.GetByUserIDAndOrganizationID(ctx, r.ID, r.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	role.Name = r.Name
	role.Permissions = r.Permissions
	role.UpdatedAt = time.Now()

	if r.GroupID != nil {
		group, err := c.groups.GetByIDAndOrganizationID(ctx, *r.GroupID, r.OrganizationID)
		if err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		role.GroupID = &group.ID
		role.Name = group.Name
		role.Permissions = group.Permissions
	}

	return c.roles.Update(ctx, role)

}
