package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvitationController struct {
	app *application.Application
}

func NewInvitationController(app *application.Application) *InvitationController {
	return &InvitationController{app: app}
}

func (c *InvitationController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)
	user := ctx.Value(models.UserContextKey).(*models.User)

	var req struct {
		GroupID         *primitive.ObjectID `json:"groupId"`
		Permissions     []string            `json:"permissions"`
		Name            string              `json:"name"`
		UserDisplayName string              `json:"userDisplayName"`
		UserEmail       string              `json:"userEmail"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request payload: " + err.Error()})
	}

	_, err = c.app.CreateInvitation(ctx, &commands.CreateInvitationRequest{
		ExpiresAt:       time.Now().Add(24 * time.Hour),
		OrganizationID:  config.OrganizationID,
		CreatorID:       user.ID,
		RoleGroupID:     req.GroupID,
		RolePermissions: req.Permissions,
		RoleName:        req.Name,
		UserDisplayName: req.UserDisplayName,
		UserEmail:       req.UserEmail,
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "invitation sent"})

}

func (c *InvitationController) Accept(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	user := r.Context().Value(models.UserContextKey).(*models.User)

	err = c.app.AcceptInvitation(r.Context(), &commands.AcceptInvitationRequest{
		InvitationID:   id,
		OrganizationID: config.OrganizationID,
		UserID:         user.ID,
		UserEmail:      user.Email,
		UserName:       user.Name,
		UserAvatarUrl:  user.AvatarUrl,
	})

	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed accepting invitation: %w", err))
		return
	}

	views.Welcome(w, r)

}
