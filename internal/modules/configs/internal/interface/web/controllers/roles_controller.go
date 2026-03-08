package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/dashboard"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/roles"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolesController struct {
	app   *application.Application
	store sessions.Store
}

func NewRolesController(app *application.Application, store sessions.Store) *RolesController {
	return &RolesController{app: app, store: store}
}

func (c *RolesController) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	var state roles.State
	var err error

	state.Roles, err = c.app.GetRoles(ctx, config.OrganizationID)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed getting roles: %w", err))
		return
	}

	roles.Page(&state).Render(ctx, w)

}

func (c *RolesController) ShowCreate(w http.ResponseWriter, r *http.Request) {
	_ = roles.Create().Render(r.Context(), w)
}

func (c *RolesController) Create(w http.ResponseWriter, r *http.Request) {

	input := struct {
		Name  string
		Email string
	}{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}

	if input.Name == "" || input.Email == "" {
		WriteHTMLErr(w, models.DecodeError(fmt.Errorf("name and email are required")))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)
	role := ctx.Value(middlewares.RoleContextKey).(*models.Role)

	request := &commands.CreateInvitationRequest{
		ExpiresAt:       time.Now().Add(24 * time.Hour),
		OrganizationID:  config.OrganizationID,
		CreatorID:       role.UserID,
		RolePermissions: []string{"*"},
		RoleName:        "Admin",
		UserDisplayName: input.Name,
		UserEmail:       input.Email,
	}

	_, err := c.app.CreateInvitation(ctx, request)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed creating invitation: %w", err))
		return
	}

	dashboard.Toast("Invitación enviada correctamente", "success").Render(ctx, w)

}

func (c *RolesController) ShowDelete(w http.ResponseWriter, r *http.Request) {

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid role id: %w", err)))
		return
	}

	name := r.URL.Query().Get("name")

	role := &models.Role{
		ID:       id,
		UserName: name,
	}

	_ = roles.Delete(role).Render(r.Context(), w)

}

func (c *RolesController) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)
	role := ctx.Value(middlewares.RoleContextKey).(*models.Role)

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid role id: %w", err)))
		return
	}

	request := &commands.DeleteRoleRequest{
		SessionRole:    role,
		TargetRoleID:   id,
		OrganizationID: config.OrganizationID,
	}

	if err := c.app.DeleteRole(ctx, request); err != nil {

		if e, ok := errors.AsType[*models.Error](err); ok {
			dashboard.Toast(e.Error(), dashboard.Danger).Render(ctx, w)
			return
		}

		slog.Error("unexpected error deleting role", "error", err)
		dashboard.Toast("Error inesperado", dashboard.Danger).Render(ctx, w)
		return
	}

	templ.Join(
		roles.DeleteRow(id),
		dashboard.Toast("El rol ha sido eliminado correctamente", "success"),
	).Render(ctx, w)

}
