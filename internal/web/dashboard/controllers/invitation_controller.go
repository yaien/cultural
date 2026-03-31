package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/application/admin"
	"github.com/yaien/cultural/internal/application/auth"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/lib/primitive"

	"github.com/yaien/cultural/internal/web/middlewares"
)

type InvitationController struct {
	invitations *admin.Invitations
}

func NewInvitationController(ivs *admin.Invitations) *InvitationController {
	return &InvitationController{ivs}
}

func (c *InvitationController) Accept(w http.ResponseWriter, r *http.Request) {
	p, err := primitive.ParseID(r.PathValue("id"))
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("invalid id: %w", err))
		return
	}

	id := primitive.ID(p)

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)
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
