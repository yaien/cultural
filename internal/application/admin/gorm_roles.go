package admin

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ RoleRepository = (*GormRoles)(nil)

type GormRoles struct {
	db *gorm.DB
}

func NewGormRoles(db *gorm.DB) *GormRoles {
	return &GormRoles{db: db}
}

func (r *GormRoles) CountAdminsByOrganizationID(ctx context.Context, organizationID primitive.ID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Where("organization_id = ? AND permissions = ?", organizationID, "*").Count(&count).Error
	return count, err
}

func (r *GormRoles) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Role, error) {
	var role Role
	err := r.db.WithContext(ctx).
		Joins("User").
		Where("id = ? AND organization_id = ?", id, organizationID).First(&role).Error
	return &role, primitive.Error(err)
}

func (r *GormRoles) GetByUserIDAndOrganizationID(ctx context.Context, userID, organizationID primitive.ID) (*Role, error) {
	var role Role
	err := r.db.WithContext(ctx).
		Joins("User").
		Where("user_id = ? AND organization_id = ?", userID, organizationID).
		First(&role).Error
	return &role, primitive.Error(err)
}

func (r *GormRoles) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) ([]*Role, error) {
	var roles []*Role
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&roles).Error
	return roles, err
}

func (r *GormRoles) Create(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *GormRoles) Update(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *GormRoles) Delete(ctx context.Context, role *Role) error {
	return r.db.WithContext(ctx).Delete(role).Error
}
