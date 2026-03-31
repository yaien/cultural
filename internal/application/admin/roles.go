package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"

	"github.com/yaien/cultural/internal/lib/coderror"
)

type Role struct {
	ID             primitive.ID `gorm:"primaryKey;autoIncrement"`
	UserID         primitive.ID `gorm:"index:idx_role_user_org,unique"`
	UserEmail      string
	UserName       string
	UserAvatarUrl  string
	OrganizationID primitive.ID `gorm:"index:idx_role_user_org,unique"`
	GroupID        *primitive.ID
	Name           string
	Permissions    Permissions
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type RoleRepository interface {
	CountAdminsByOrganizationID(ctx context.Context, organizationID primitive.ID) (int64, error)
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Role, error)
	GetByUserIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Role, error)
	GetByOrganizationID(ctx context.Context, id primitive.ID) ([]*Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, role *Role) error
}

type Roles struct {
	repository RoleRepository
}

func NewRoles(repository RoleRepository) *Roles {
	return &Roles{repository: repository}
}

func (q *Roles) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) ([]*Role, error) {
	return q.repository.GetByOrganizationID(ctx, organizationID)
}

func (q *Roles) GetByUserIDAndOrganizationID(ctx context.Context, userID, organizationID primitive.ID) (*Role, error) {
	return q.repository.GetByUserIDAndOrganizationID(ctx, userID, organizationID)
}

func (c *Roles) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Role, error) {
	return c.repository.GetByIDAndOrganizationID(ctx, id, organizationID)
}

func (c *Roles) Create(ctx context.Context, role *Role) error {
	return c.repository.Create(ctx, role)
}

type DeleteRoleOptions struct {
	SessionRole    *Role
	TargetRoleID   primitive.ID
	OrganizationID primitive.ID
}

func (c *Roles) Delete(ctx context.Context, req *DeleteRoleOptions) error {
	if req.SessionRole == nil {
		return fmt.Errorf("session role is required")
	}

	if req.SessionRole.ID == req.TargetRoleID {
		return coderror.Newf("invalid_deletion", "no puedes borrar tu propio rol")
	}

	if !req.SessionRole.Permissions.Has("delete-role") {
		return coderror.Newf("permission_denied", "no tienes permiso para eliminar roles")
	}

	count, err := c.repository.CountAdminsByOrganizationID(ctx, req.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to count roles: %w", err)
	}

	if count == 1 {
		return coderror.Newf("last_admin_deletion", "no puedes eliminar el último rol de administrador")
	}

	role, err := c.repository.GetByIDAndOrganizationID(ctx, req.TargetRoleID, req.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	err = c.repository.Delete(ctx, role)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}
