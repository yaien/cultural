package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/integrations"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type IntegrationController struct {
	app      *application.Application
	registry *models.IntegrationRegistry
}

func NewIntegrationController(app *application.Application, registry *models.IntegrationRegistry) *IntegrationController {
	return &IntegrationController{app: app, registry: registry}
}

func (c *IntegrationController) Index(w http.ResponseWriter, r *http.Request) {
	integrations.Page(c.registry.All()).Render(r.Context(), w)
}

func (c *IntegrationController) Integration(w http.ResponseWriter, r *http.Request) {
	def, ok := c.registry.Get(r.PathValue("integration"))
	if !ok {
		WriteHTMLErr(w, models.NotFoundError(fmt.Errorf("itegration %q not found", r.PathValue("integration"))))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	page, err := def.Page(ctx, config)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed getting integration page: %w", err))
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		_ = integrations.Modal(def, page).Render(ctx, w)
		return
	}

	definitions := c.registry.All()
	_ = integrations.Detail(definitions, def, page).Render(ctx, w)
}

func (c *IntegrationController) OAuthLogin(w http.ResponseWriter, r *http.Request) {

	def, ok := c.registry.Get(r.PathValue("integration"))
	if !ok {
		WriteHTMLErr(w, models.NotFoundError(fmt.Errorf("itegration %q not found", r.PathValue("integration"))))
		return
	}

	lgn, ok := def.(models.IntegrationOAuth)
	if !ok {
		WriteHTMLErr(w, models.NotFoundError(fmt.Errorf("integration %q does not support oauth", r.PathValue("integration"))))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	url, err := lgn.OAuthCodeURL(ctx, config)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed getting oauth code url: %w", err))
		return
	}

	http.Redirect(w, r, url, http.StatusPermanentRedirect)

}

func (c *IntegrationController) OauthCallback(w http.ResponseWriter, r *http.Request) {

	def, ok := c.registry.Get(r.PathValue("integration"))
	if !ok {
		WriteHTMLErr(w, models.NotFoundError(fmt.Errorf("itegration %q not found", r.PathValue("integration"))))
		return
	}

	lgn, ok := def.(models.IntegrationOAuth)
	if !ok {
		WriteHTMLErr(w, models.NotFoundError(fmt.Errorf("integration %q does not support oauth", r.PathValue("integration"))))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	code := r.URL.Query().Get("code")
	err := lgn.OAuthExchange(ctx, config, code)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed exchanging code for token: %w", err))
		return
	}

	http.Redirect(w, r, "/dashboard/integrations/"+def.Name(), http.StatusPermanentRedirect)
}
