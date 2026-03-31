package integration

import (
	"context"
	"text/template"
	"time"

	"github.com/a-h/templ"
	"github.com/robfig/cron/v3"
	"github.com/yaien/cultural/internal/application/label"

	"github.com/yaien/cultural/internal/lib/primitive"
	"github.com/yaien/cultural/internal/lib/worker"
)

type Integration[T any] struct {
	ID             primitive.ID `gorm:"primaryKey,autoIncrement"`
	OrganizationID primitive.ID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Data           T `gorm:"type:jsonb;serializer:json"`
}

func (i *Integration[T]) TableName() string {
	return "integrations"
}

type GetOptions struct {
	OrganizationID primitive.ID
	Name           string
}

type Repository[T any] interface {
	Create(ctx context.Context, i *Integration[T]) error
	Update(ctx context.Context, i *Integration[T]) error
	GetByOrganizationIDAndName(ctx context.Context, id primitive.ID, name string) (*Integration[T], error)
	GetByName(ctx context.Context, name string) ([]*Integration[T], error)
}

type Definition interface {
	Name() string
	Title() string
	Description() string
	Image() string
	Page(ctx context.Context, config *label.Config) (templ.Component, error)
}

type OAuth interface {
	OAuthCodeURL(ctx context.Context, config *label.Config) (url string, err error)
	OAuthExchange(ctx context.Context, config *label.Config, code string) error
}

type TemplateFuncMap interface {
	TemplateFuncMap(ctx context.Context, config *label.Config) template.FuncMap
}

type Background interface {
	RegisterBackgroundProcess(cron *cron.Cron, queue *worker.Queue, wk *worker.Worker)
}

type Registry struct {
	definitions map[string]Definition
}

func NewRegistry(def ...Definition) *Registry {
	d := &Registry{
		definitions: make(map[string]Definition),
	}

	for _, def := range def {
		d.Register(def)
	}

	return d
}

func (r *Registry) Register(d Definition) {
	r.definitions[d.Name()] = d
}

func (r *Registry) Get(name string) (Definition, bool) {
	d, ok := r.definitions[name]
	return d, ok
}

func (r *Registry) All() []Definition {
	definitions := make([]Definition, 0, len(r.definitions))
	for _, d := range r.definitions {
		definitions = append(definitions, d)
	}
	return definitions
}
