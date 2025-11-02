package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetRoleQuery struct {
	roles models.RoleRepository
}

func NewGetRoleQuery(roles models.RoleRepository) *GetRoleQuery {
	return &GetRoleQuery{roles: roles}
}

func (q *GetRoleQuery) GetRole(ctx context.Context, userId, organizationId primitive.ObjectID) (*models.Role, error) {
	return q.roles.GetByUserIDAndOrganizationID(ctx, userId, organizationId)
}
