package models

import (
	"context"
	"text/template"
	"time"

	"github.com/a-h/templ"
	"github.com/robfig/cron/v3"
	"github.com/yaien/cultural/internal/library/worker"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Integration[T any] struct {
	ID             primitive.ObjectID `bson:"_id"`
	OrganizationID primitive.ObjectID `bson:"organizationId"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	Name           string             `bson:"name"`
	Data           T                  `bson:"data"`
}

type GetIntegrationOptions struct {
	OrganizationID primitive.ObjectID
	Name           string
}

type IntegrationRepository[T any] interface {
	Create(ctx context.Context, i *Integration[T]) error
	Update(ctx context.Context, i *Integration[T]) error
	GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) (*Integration[T], error)
	GetByName(ctx context.Context, name string) ([]*Integration[T], error)
}

type IntegrationDefinition interface {
	Name() string
	Title() string
	Description() string
	Image() string
	Page(ctx context.Context, config *Config) (templ.Component, error)
}

type IntegrationOAuth interface {
	OAuthCodeURL(ctx context.Context, config *Config) (url string, err error)
	OAuthExchange(ctx context.Context, config *Config, code string) error
}

type IntegrationTemplateFuncMap interface {
	TemplateFuncMap(ctx context.Context, config *Config) template.FuncMap
}

type IntegrationBackground interface {
	RegisterBackgroundProcess(cron *cron.Cron, queue *worker.Queue, wk *worker.Worker)
}

type IntegrationRegistry struct {
	definitions map[string]IntegrationDefinition
}

func NewIntegrationRegistry(def ...IntegrationDefinition) *IntegrationRegistry {
	d := &IntegrationRegistry{
		definitions: make(map[string]IntegrationDefinition),
	}

	for _, def := range def {
		d.Register(def)
	}

	return d
}

func (r *IntegrationRegistry) Register(d IntegrationDefinition) {
	r.definitions[d.Name()] = d
}

func (r *IntegrationRegistry) Get(name string) (IntegrationDefinition, bool) {
	d, ok := r.definitions[name]
	return d, ok
}

func (r *IntegrationRegistry) All() []IntegrationDefinition {
	definitions := make([]IntegrationDefinition, 0, len(r.definitions))
	for _, d := range r.definitions {
		definitions = append(definitions, d)
	}
	return definitions
}
