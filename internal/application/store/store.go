package store

import (
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"

	"github.com/yaien/cultural/internal/application/storage"
)

type Product struct {
	ID             primitive.ID `gorm:"primaryKey;autoIncrement"`
	OrganizationID primitive.ID `gorm:"index"`
	Slug           string       `gorm:"index"`
	Name           string
	Presentations  []Presentation `gorm:"type:jsonb;serializer:json"`
	Published      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Presentation struct {
	ID       primitive.UUID
	Contents []Content
	Name     string
	Quantity int
	Price    float64
}

type Content struct {
	ID         primitive.UUID
	FileID     primitive.ID
	FilePreset string
}

type Store struct {
	Products      *Products
	Presentations *Presentations
	Contents      *Contents
}

func New(db *gorm.DB, storage *storage.Storage) *Store {
	return &Store{
		Products:      NewProducts(db),
		Presentations: NewPresentations(db, storage),
		Contents:      NewContents(db, storage),
	}
}
