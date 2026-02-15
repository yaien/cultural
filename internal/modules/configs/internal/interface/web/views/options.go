package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type data struct {
	Config *models.Config
	User   *models.User
	Role   *models.Role
	Path   string
	Meta   any
}

type option func(*data)

func Meta(meta any) option {
	return func(d *data) {
		d.Meta = meta
	}
}

func newData(r *http.Request, opts ...option) *data {
	var data data

	data.Path = r.URL.Path

	if config, ok := r.Context().Value(models.ConfigContextKey).(*models.Config); ok {
		data.Config = config
	}

	if user, ok := r.Context().Value(models.UserContextKey).(*models.User); ok {
		data.User = user
	}

	if role, ok := r.Context().Value(models.RoleContextKey).(*models.Role); ok {
		data.Role = role
	}

	for _, opt := range opts {
		opt(&data)
	}

	return &data
}

type Link struct {
	Path   string
	Name   string
	Icon   template.HTML
	Active bool
}

func (d *data) Link(path, name, icon string) (*Link, error) {
	svg, err := fs.ReadFile(fmt.Sprintf("icons/%s", icon))
	if err != nil {
		return nil, fmt.Errorf("failed reading icon: %w", err)
	}

	return &Link{
		Path:   path,
		Name:   name,
		Icon:   template.HTML(svg),
		Active: path == d.Path,
	}, nil
}
