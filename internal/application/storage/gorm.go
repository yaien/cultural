package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ Repository = (*Gorm)(nil)

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *Gorm {
	return &Gorm{db: db}
}

func (r *Gorm) Create(ctx context.Context, file *File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *Gorm) Update(ctx context.Context, file *File) error {
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *Gorm) GetByID(ctx context.Context, id primitive.ID) (*File, error) {
	var file File
	err := r.db.WithContext(ctx).First(&file, id).Error
	switch {
	case err == nil:
		return &file, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New("not_found", err)
	default:
		return nil, fmt.Errorf("failed finding file: %w", err)
	}
}

func (r *Gorm) GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ID, name string) (*File, error) {
	var file File
	err := r.db.WithContext(ctx).Where("organization_id = ? AND name = ?", organizationID, name).First(&file).Error
	switch {
	case err == nil:
		return &file, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New("not_found", err)
	default:
		return nil, fmt.Errorf("failed finding file: %w", err)
	}
}

func (r *Gorm) GetByOrganizationIDAndID(ctx context.Context, organizationID primitive.ID, id primitive.ID) (*File, error) {
	var file File
	err := r.db.WithContext(ctx).Where("organization_id = ? AND id = ?", organizationID, id).First(&file).Error
	switch {
	case err == nil:
		return &file, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New("not_found", err)
	default:
		return nil, fmt.Errorf("failed finding file: %w", err)
	}
}

func (r *Gorm) DeleteByOrganizationIDAndName(ctx context.Context, organizationID primitive.ID, name string) error {
	return r.db.
		WithContext(ctx).
		Where("organization_id = ? AND name = ?", organizationID, name).
		Delete(&File{}).Error
}

func (r *Gorm) RenameByOrganizationIDAndName(ctx context.Context, organizationID primitive.ID, oldName string, newName string) error {
	return r.db.WithContext(ctx).
		Model(&File{}).
		Where("organization_id = ? AND name = ?", organizationID, oldName).
		Update("name", newName).Error
}

func (r *Gorm) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) (files []*File, err error) {
	err = r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&files).Error
	return
}
