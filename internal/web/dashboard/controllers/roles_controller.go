package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/application/admin"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/web/dashboard/views/dashboard"
	"github.com/yaien/cultural/internal/web/dashboard/views/roles"
	"github.com/yaien/cultural/internal/web/middlewares"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolesController struct {
	roles       *admin.Roles
	invitations *admin.Invitations
	store       sessions.Store
}

func NewRolesController(rls *admin.Roles, ivs *admin.Invitations, store sessions.Store) *RolesController {
	return &RolesController{rls, ivs, store}
}

func (c *RolesController) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	var state roles.State
	var err error

	state.Roles, err = c.roles.GetByOrganizationID(ctx, config.OrganizationID)
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
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "name and email are required"))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)
	role := ctx.Value(middlewares.RoleContextKey).(*admin.Role)

	request := &admin.CreateInvitationOptions{
		ExpiresAt:       time.Now().Add(24 * time.Hour),
		OrganizationID:  config.OrganizationID,
		CreatorID:       role.UserID,
		RolePermissions: []string{"*"},
		RoleName:        "Admin",
		UserDisplayName: input.Name,
		UserEmail:       input.Email,
	}

	_, err := c.invitations.Create(ctx, request)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed creating invitation: %w", err))
		return
	}

	dashboard.Toast("Invitación enviada correctamente", "success").Render(ctx, w)

}

func (c *RolesController) ShowDelete(w http.ResponseWriter, r *http.Request) {

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteJSONErr(w, coderror.Newf(coderror.DecodeFailed, "invalid role id: %w", err))
		return
	}

	name := r.URL.Query().Get("name")

	role := &admin.Role{
		ID:       id,
		UserName: name,
	}

	_ = roles.Delete(role).Render(r.Context(), w)

}

func (c *RolesController) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)
	role := ctx.Value(middlewares.RoleContextKey).(*admin.Role)

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteJSONErr(w, coderror.Newf(coderror.DecodeFailed, "invalid role id: %w", err))
		return
	}

	request := &admin.DeleteRoleOptions{
		SessionRole:    role,
		TargetRoleID:   id,
		OrganizationID: config.OrganizationID,
	}

	if err := c.roles.Delete(ctx, request); err != nil {

		if e, ok := errors.AsType[*coderror.Error](err); ok {
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
