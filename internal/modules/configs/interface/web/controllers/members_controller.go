package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/views"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MembersController struct {
	app *application.Application
}

func NewMembersController(app *application.Application) *MembersController {
	return &MembersController{app: app}
}

func (c *MembersController) Index(w http.ResponseWriter, r *http.Request) {
	views.Roles(w, r)
}

func (c *MembersController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	roles, err := c.app.GetRoles(ctx, config.OrganizationID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)

}

func (c *MembersController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid role id: " + err.Error()})
		return
	}

	var input struct {
		GroupID     *primitive.ObjectID `json:"groupId"`
		Permissions []string            `json:"permissions"`
		Name        string              `json:"name"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request payload: " + err.Error()})
		return
	}

	err = c.app.UpdateRole(ctx, &commands.UpdateRoleRequest{
		ID:             id,
		OrganizationID: config.OrganizationID,
		GroupID:        input.GroupID,
		Permissions:    input.Permissions,
		Name:           input.Name,
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "role updated"})
}

func (c *MembersController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid role id: " + err.Error()})
		return
	}

	err = c.app.DeleteRole(ctx, id, config.OrganizationID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "role deleted"})
}
