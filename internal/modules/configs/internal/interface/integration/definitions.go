package integration

import (
	"context"

	"github.com/a-h/templ"
	"github.com/markbates/goth"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

//go:generate go tool templ generate

type Definition struct {
	Name        string
	Description string
	Image       string
	Path        string
	Data        any
	HandleOauth func(ctx context.Context, oauth *OAuth) (err error)
	PageSection func(integration *models.Integration) templ.Component
}

type OAuth struct {
	Config *models.Config
	App    *application.Application
	User   *goth.User
}

var Definitions = map[string]*Definition{
	Instagram.Path: Instagram,
}
