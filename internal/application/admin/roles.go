package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/lib/coderror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty"`
	UserID         primitive.ObjectID  `bson:"userId"`
	UserEmail      string              `bson:"userEmail"`
	UserName       string              `bson:"userName"`
	UserAvatarUrl  string              `bson:"userAvatarUrl"`
	OrganizationID primitive.ObjectID  `bson:"organizationId"`
	GroupID        *primitive.ObjectID `bson:"groupId,omitempty"`
	Name           string              `bson:"name"`
	Permissions    Permissions         `bson:"permissions"`
	CreatedAt      time.Time           `bson:"createdAt"`
	UpdatedAt      time.Time           `bson:"updatedAt"`
	DeletedAt      *time.Time          `bson:"deletedAt,omitempty"`
}

type RoleRepository interface {
	CountAdminsByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (int64, error)
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Role, error)
	GetByUserIDAndOrganizationID(ctx context.Context, userId, organizationID primitive.ObjectID) (*Role, error)
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*Role, error)
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

func (q *Roles) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*Role, error) {
	return q.repository.GetByOrganizationID(ctx, organizationID)
}

func (q *Roles) GetByUserIDAndOrganizationID(ctx context.Context, userId, organizationID primitive.ObjectID) (*Role, error) {
	return q.repository.GetByUserIDAndOrganizationID(ctx, userId, organizationID)
}

func (c *Roles) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Role, error) {
	return c.repository.GetByIDAndOrganizationID(ctx, id, organizationID)
}

func (c *Roles) Create(ctx context.Context, role *Role) error {
	return c.repository.Create(ctx, role)
}

type DeleteRoleOptions struct {
	SessionRole    *Role
	TargetRoleID   primitive.ObjectID
	OrganizationID primitive.ObjectID
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
