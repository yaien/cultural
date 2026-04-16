package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/application/auth"
	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"

	"github.com/yaien/cultural/internal/lib/coderror"
)

type Role struct {
	ID             primitive.ID `gorm:"primaryKey;autoIncrement"`
	UserID         primitive.ID `gorm:"index:idx_role_user_org,unique"`
	User           *auth.User
	OrganizationID primitive.ID `gorm:"index:idx_role_user_org,unique"`
	GroupID        *primitive.ID
	Name           string
	Permissions    Permissions `gorm:"type:jsonb;serializer:json"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type Roles struct {
	repository gorm.Interface[Role]
}

func NewRoles(db *gorm.DB) *Roles {
	return &Roles{gorm.G[Role](db)}
}

func (q *Roles) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) ([]Role, error) {
	return q.repository.
		Preload("User", nil).
		Where("organization_id = ?", organizationID).Find(ctx)
}

func (q *Roles) GetByUserIDAndOrganizationID(ctx context.Context, userID, organizationID primitive.ID) (Role, error) {
	return q.repository.
		Preload("User", nil).
		Where("user_id = ? AND organization_id = ?", userID, organizationID).First(ctx)
}

func (c *Roles) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (Role, error) {
	return c.repository.
		Preload("User", nil).
		Where("id = ? AND organization_id = ?", id, organizationID).First(ctx)
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

	count, err := c.repository.
		Where("organization_id = ? AND permissions = ?", req.OrganizationID, "*").
		Count(ctx, "*")

	if err != nil {
		return fmt.Errorf("failed to count roles: %w", err)
	}

	if count == 1 {
		return coderror.Newf("last_admin_deletion", "no puedes eliminar el último rol de administrador")
	}

	deleted, err := c.repository.Where("id = ? AND organization_id = ?", req.TargetRoleID, req.OrganizationID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if deleted == 0 {
		return coderror.Newf("not_found", "no se encontró el rol a eliminar")
	}

	return nil
}
