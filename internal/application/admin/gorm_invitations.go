package admin

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ InvitationRepository = (*GormInvitations)(nil)

type GormInvitations struct {
	db *gorm.DB
}

func NewGormInvitations(db *gorm.DB) *GormInvitations {
	return &GormInvitations{db: db}
}

func (r *GormInvitations) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Invitation, error) {
	var invitation Invitation
	err := r.db.WithContext(ctx).Where("id = ? AND organization_id = ?", id, organizationID).First(&invitation).Error
	return &invitation, primitive.Error(err)
}

func (r *GormInvitations) Create(ctx context.Context, invitation *Invitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *GormInvitations) Update(ctx context.Context, invitation *Invitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}
