package controllers

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/markbates/goth/gothic"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/integration"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/integrations"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type IntegrationController struct {
	app *application.Application
}

func NewIntegrationController(app *application.Application) *IntegrationController {
	return &IntegrationController{app: app}
}

func (c *IntegrationController) Index(w http.ResponseWriter, r *http.Request) {
	integrations.Page().Render(r.Context(), w)
}

func (c *IntegrationController) Integration(w http.ResponseWriter, r *http.Request) {
	def, ok := integration.Definitions[r.PathValue("integration")]
	if !ok {
		WriteHTMLErr(w, models.NotFoundError(fmt.Errorf("itegration %q not found", r.PathValue("integration"))))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	intgr, err := c.app.GetIntegration(ctx, queries.GetIntegrationOptions{
		OrganizationID: config.OrganizationID,
		Name:           def.Name,
		Data:           reflect.New(reflect.TypeOf(def.Data)).Elem().Interface(),
	})

	if err != nil && !models.IsNotFoundError(err) {
		WriteHTMLErr(w, fmt.Errorf("failed getting integration: %w", err))
		return
	}

	integrations.Integration(def, intgr).Render(ctx, w)

}

func (c *IntegrationController) OauthConnect(w http.ResponseWriter, r *http.Request) {

	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}

	definition, ok := integration.Definitions[u.Provider]
	if !ok {
		WriteHTMLErr(w, fmt.Errorf("integration not found for provider %s: %w", u.Provider, err))
		return
	}

	if definition.HandleOauth == nil {
		WriteHTMLErr(w, fmt.Errorf("integration %s does not support oauth: %w", u.Provider, err))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	oauth := &integration.OAuth{
		Config: config,
		App:    c.app,
		User:   &u,
	}

	if err = definition.HandleOauth(ctx, oauth); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed handling oauth for integration: %w", err))
		return
	}

	http.Redirect(w, r, "/dashboard/integration/"+definition.Name, http.StatusPermanentRedirect)

}

func (c *IntegrationController) OauthCallback(w http.ResponseWriter, r *http.Request) {

	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return
	}

	definition, ok := integration.Definitions[u.Provider]
	if !ok {
		WriteHTMLErr(w, fmt.Errorf("integration not found for provider %s: %w", u.Provider, err))
		return
	}

	if definition.HandleOauth == nil {
		WriteHTMLErr(w, fmt.Errorf("integration %s does not support oauth: %w", u.Provider, err))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	oauth := &integration.OAuth{
		Config: config,
		App:    c.app,
		User:   &u,
	}

	if err = definition.HandleOauth(ctx, oauth); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed handling oauth for integration: %w", err))
		return
	}

	http.Redirect(w, r, "/dashboard/integration/"+definition.Name, http.StatusPermanentRedirect)
}
