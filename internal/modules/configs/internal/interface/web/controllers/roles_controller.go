package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolesController struct {
	app *application.Application
}

func NewRolesController(app *application.Application) *RolesController {
	return &RolesController{app: app}
}

func (c *RolesController) Index(w http.ResponseWriter, r *http.Request) {
	views.Roles(w, r)
}

func (c *RolesController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	roles, err := c.app.GetRoles(ctx, config.OrganizationID)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed getting roles: %w", err))
		return
	}

	WriteJSON(w, roles)

}

func (c *RolesController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid role id: %w", err)))
		return
	}

	var input struct {
		GroupID     *primitive.ObjectID `json:"groupId"`
		Permissions []string            `json:"permissions"`
		Name        string              `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid request payload: %w", err)))
		return
	}

	request := &commands.UpdateRoleRequest{
		ID:             id,
		OrganizationID: config.OrganizationID,
		GroupID:        input.GroupID,
		Permissions:    input.Permissions,
		Name:           input.Name,
	}

	if err := c.app.UpdateRole(ctx, request); err != nil {
		WriteJSONErr(w, fmt.Errorf("failed updating role: %w", err))
		return
	}

	WriteJSONSuccess(w)
}

func (c *RolesController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)
	sessionRole := ctx.Value(models.RoleContextKey).(*models.Role)

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid role id: %w", err)))
		return
	}

	request := &commands.DeleteRoleRequest{
		SessionRole:    sessionRole,
		TargetRoleID:   id,
		OrganizationID: config.OrganizationID,
	}

	if err := c.app.DeleteRole(ctx, request); err != nil {
		WriteJSONErr(w, fmt.Errorf("failed deleting role: %w", err))
		return
	}

}
