package integration

import (
	"context"
	"errors"

	"github.com/yaien/cultural/internal/lib/primitive"

	"github.com/yaien/cultural/internal/lib/coderror"
	"gorm.io/gorm"
)

var _ Repository[any] = (*Gorm[any])(nil)

type Gorm[T any] struct {
	db *gorm.DB
}

func NewGorm[T any](db *gorm.DB) *Gorm[T] {
	return &Gorm[T]{db: db}
}

func (r *Gorm[T]) Create(ctx context.Context, i *Integration[T]) error {
	return r.db.WithContext(ctx).Create(i).Error
}

func (r *Gorm[T]) Update(ctx context.Context, i *Integration[T]) error {
	return r.db.WithContext(ctx).Save(i).Error
}

func (r *Gorm[T]) GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ID, name string) (*Integration[T], error) {
	var i Integration[T]
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND name = ?", organizationID, name).
		First(&i).Error
	switch {
	case err == nil:
		return &i, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *Gorm[T]) GetByName(ctx context.Context, name string) ([]*Integration[T], error) {
	var integrations []*Integration[T]
	err := r.db.WithContext(ctx).Where("name = ?", name).Find(&integrations).Error
	if err != nil {
		return nil, err
	}
	return integrations, nil
}
