package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/library/admin"
	"github.com/yaien/cultural/internal/library/auth"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvitationController struct {
	invitations *admin.Invitations
}

func NewInvitationController(ivs *admin.Invitations) *InvitationController {
	return &InvitationController{ivs}
}

func (c *InvitationController) Accept(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)
	user := ctx.Value(middlewares.UserContextKey).(*auth.User)

	opts := &admin.AcceptInvitationOptions{
		InvitationID:   id,
		OrganizationID: config.OrganizationID,
		UserID:         user.ID,
		UserEmail:      user.Email,
		UserName:       user.Name,
		UserAvatarUrl:  user.AvatarUrl,
	}

	if err := c.invitations.Accept(ctx, opts); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed accepting invitation: %w", err))
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)

}
