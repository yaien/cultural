package admin

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ RoleRepository = (*GormRoles)(nil)

type GormRoles struct {
	db gorm.Interface[Role]
}

func NewGormRoles(db *gorm.DB) *GormRoles {
	return &GormRoles{db: gorm.G[Role](db)}
}

func (r *GormRoles) CountAdminsByOrganizationID(ctx context.Context, organizationID primitive.ID) (int64, error) {
	return r.db.Where("organization_id = ? AND permissions = ?", organizationID, "*").Count(ctx, "*")
}

func (r *GormRoles) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Role, error) {
	role, err := r.db.
		Joins(clause.JoinTarget{Association: "User"}, nil).
		Where("id = ? AND organization_id = ?", id, organizationID).
		First(ctx)
	return &role, primitive.Error(err)
}

func (r *GormRoles) GetByUserIDAndOrganizationID(ctx context.Context, userID, organizationID primitive.ID) (*Role, error) {
	role, err := r.db.
		Joins(clause.JoinTarget{Association: "User"}, nil).
		Where("user_id = ? AND organization_id = ?", userID, organizationID).
		First(ctx)
	return &role, primitive.Error(err)
}

func (r *GormRoles) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) ([]Role, error) {
	roles, err := r.db.
		Preload("User", nil).
		Where("organization_id = ?", organizationID).
		Find(ctx)
	return roles, err
}

func (r *GormRoles) Create(ctx context.Context, role *Role) error {
	return r.db.Create(ctx, role)
}

func (r *GormRoles) Update(ctx context.Context, role *Role) error {
	_, err := r.db.Updates(ctx, *role)
	return err
}

func (r *GormRoles) Delete(ctx context.Context, role *Role) error {
	_, err := r.db.Where("id = ?", role.ID).Delete(ctx)
	return err
}
