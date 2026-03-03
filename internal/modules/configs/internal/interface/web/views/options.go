package views

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type data struct {
	Config   *models.Config
	User     *models.User
	Role     *models.Role
	Path     string
	Data     any
	Template string
}

type option func(*data)

func Data(dat any) option {
	return func(d *data) {
		d.Data = dat
	}
}

func Template(name string) option {
	return func(d *data) {
		d.Template = name
	}
}

func newData(r *http.Request, opts ...option) *data {
	var data data

	data.Path = r.URL.Path

	if config, ok := r.Context().Value(middlewares.ConfigContextKey).(*models.Config); ok {
		data.Config = config
	}

	if user, ok := r.Context().Value(middlewares.UserContextKey).(*models.User); ok {
		data.User = user
	}

	if role, ok := r.Context().Value(middlewares.RoleContextKey).(*models.Role); ok {
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
	Icon   string
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
		Icon:   string(svg),
		Active: path == d.Path,
	}, nil
}
