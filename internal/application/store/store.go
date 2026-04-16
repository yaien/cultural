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
	Presentations  []*Presentation
	Published      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Presentation struct {
	ID        primitive.ID `gorm:"primaryKey;autoIncrement"`
	ProductID primitive.ID
	Contents  []*Content
	Name      string
	Quantity  int
	Price     float64
}

type Content struct {
	ID             primitive.ID `gorm:"primaryKey;autoIncrement"`
	PresentationID primitive.ID
	FileID         primitive.ID
	File           storage.File
	Order          int
}

type Store struct {
	Products      *Products
	Presentations *Presentations
	Files         *Files
}

func New(db *gorm.DB, storage *storage.Storage) *Store {
	return &Store{
		Products:      NewProducts(db),
		Presentations: NewPresentations(db),
		Files:         NewFiles(db, storage),
	}
}
