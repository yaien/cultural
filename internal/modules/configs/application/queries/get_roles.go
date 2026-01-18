package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetRolesQuery struct {
	roles models.RoleRepository
}

func NewGetRolesQuery(roles models.RoleRepository) *GetRolesQuery {
	return &GetRolesQuery{
		roles: roles,
	}
}

func (q *GetRolesQuery) GetRoles(ctx context.Context, organizationID primitive.ObjectID) ([]*models.Role, error) {
	return q.roles.GetByOrganizationID(ctx, organizationID)
}
