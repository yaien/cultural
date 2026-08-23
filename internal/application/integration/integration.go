package integration

import (
	"context"
	"maps"
	"net/http"
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

type Definition interface {
	Name() string
}

type Page interface {
	PageTitle() string
	PageDescription() string
	PageImage() string
	PageComponent(ctx context.Context, config *label.Config) (templ.Component, error)
}

type OAuth interface {
	OAuthCodeURL(ctx context.Context, config *label.Config) (url string, err error)
	OAuthExchange(ctx context.Context, config *label.Config, code string) error
}

type TemplatePreset struct {
	Title   string
	Options func() ([]TemplateOption, error)
	Load    func(params ...string) (any, error)
}

type TemplateOption struct {
	Value string
	Label string
}

type TemplateActionCall struct {
	Headers         map[string]string
	Body            string
	ResponseWritter http.ResponseWriter
	Request         *http.Request
}

type TemplateAction struct {
	Title  string
	Handle func(r *http.Request) (any, error)
}

type TemplateFuncMap = template.FuncMap
type TemplatePresetMap = map[string]TemplatePreset
type TemplateActionMap = map[string]TemplateAction

type TemplateFuncMapper interface {
	TemplateFuncMap(ctx context.Context, config *label.Config) TemplateFuncMap
}

type TemplatePresetMapper interface {
	TemplatePresetMap(ctx context.Context, config *label.Config) TemplatePresetMap
}

type TemplateActionMapper interface {
	TemplateActionMap(ctx context.Context, config *label.Config) TemplateActionMap
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

func (r *Registry) TemplateFuncMap(ctx context.Context, config *label.Config) map[string]any {
	funcs := make(TemplateFuncMap)
	for _, d := range r.definitions {
		if m, ok := d.(TemplateFuncMapper); ok {
			maps.Copy(funcs, m.TemplateFuncMap(ctx, config))
		}
	}
	return funcs
}

func (r *Registry) TemplatePresetMap(ctx context.Context, config *label.Config) TemplatePresetMap {
	presets := make(TemplatePresetMap)
	for _, d := range r.definitions {
		if m, ok := d.(TemplatePresetMapper); ok {
			maps.Copy(presets, m.TemplatePresetMap(ctx, config))
		}
	}
	return presets
}

func (r *Registry) TemplateActionMap(ctx context.Context, config *label.Config) TemplateActionMap {
	actions := make(TemplateActionMap)
	for _, d := range r.definitions {
		if m, ok := d.(TemplateActionMapper); ok {
			maps.Copy(actions, m.TemplateActionMap(ctx, config))
		}
	}
	return actions
}
