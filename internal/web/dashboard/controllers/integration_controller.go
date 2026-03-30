package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/web/dashboard/views/integrations"
	"github.com/yaien/cultural/internal/web/middlewares"
)

type IntegrationController struct {
	registry *integration.Registry
}

func NewIntegrationController(registry *integration.Registry) *IntegrationController {
	return &IntegrationController{registry: registry}
}

func (c *IntegrationController) Index(w http.ResponseWriter, r *http.Request) {
	_ = integrations.Page(c.registry.All()).Render(r.Context(), w)
}

func (c *IntegrationController) Integration(w http.ResponseWriter, r *http.Request) {
	def, ok := c.registry.Get(r.PathValue("integration"))
	if !ok {
		WriteHTMLErr(w, coderror.Newf(coderror.NotFound, "integration not found"))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

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
		WriteHTMLErr(w, coderror.Newf(coderror.NotFound, "integration %q not found", r.PathValue("integration")))
		return
	}

	lgn, ok := def.(integration.OAuth)
	if !ok {
		WriteHTMLErr(w, coderror.Newf(coderror.NotFound, "integration %q does not support oauth", r.PathValue("integration")))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

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
		WriteHTMLErr(w, coderror.Newf(coderror.NotFound, "integration %q not found", r.PathValue("integration")))
		return
	}

	lgn, ok := def.(integration.OAuth)
	if !ok {
		WriteHTMLErr(w, coderror.Newf(coderror.NotFound, "integration %q does not support oauth", r.PathValue("integration")))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	code := r.URL.Query().Get("code")
	err := lgn.OAuthExchange(ctx, config, code)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed exchanging code for token: %w", err))
		return
	}

	http.Redirect(w, r, "/dashboard/integrations/"+def.Name(), http.StatusFound)
}
