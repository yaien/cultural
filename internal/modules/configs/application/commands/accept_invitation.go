package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AcceptInvitationCommand struct {
	invitations models.InvitationRepository
	roles       models.RoleRepository
}

func NewAcceptInvitationCommand(invitations models.InvitationRepository, roles models.RoleRepository) *AcceptInvitationCommand {
	return &AcceptInvitationCommand{
		invitations: invitations,
		roles:       roles,
	}
}

type AcceptInvitationRequest struct {
	InvitationID   primitive.ObjectID
	OrganizationID primitive.ObjectID
	UserID         primitive.ObjectID
	UserEmail      string
	UserName       string
	UserAvatarUrl  string
}

func (c *AcceptInvitationCommand) AcceptInvitation(ctx context.Context, req *AcceptInvitationRequest) error {
	invitation, err := c.invitations.GetByIDAndOrganizationID(ctx, req.InvitationID, req.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed getting invitation: %w", err)
	}

	if invitation.UserEmail != req.UserEmail {
		return &models.Error{Code: "email_missmatch", Err: fmt.Errorf("%s user email doesn't match with this invitation", req.UserEmail)}
	}

	if invitation.AcceptedAt != nil {
		return &models.Error{Code: "invitation_already_accepted", Err: fmt.Errorf("invitation %s already accepted", req.InvitationID.Hex())}
	}

	now := time.Now()
	invitation.AcceptedAt = &now

	err = c.invitations.Update(ctx, invitation)
	if err != nil {
		return fmt.Errorf("failed updating invitation: %w", err)
	}

	role, err := c.roles.GetByUserIDAndOrganizationID(ctx, req.UserID, req.OrganizationID)

	var e *models.Error

	switch {
	case err == nil:
		role.Permissions = invitation.RolePermissions
		role.Name = invitation.RoleName
		role.GroupID = invitation.RoleGroupID
		role.UpdatedAt = time.Now()

		err = c.roles.Update(ctx, role)
		if err != nil {
			return fmt.Errorf("failed updating role: %w", err)
		}

		return nil

	case errors.As(err, &e) && e.Code == "not_found":
		role = &models.Role{
			UserID:         req.UserID,
			UserEmail:      req.UserEmail,
			UserName:       req.UserName,
			UserAvatarUrl:  req.UserAvatarUrl,
			OrganizationID: req.OrganizationID,
			Name:           invitation.RoleName,
			Permissions:    invitation.RolePermissions,
			GroupID:        invitation.RoleGroupID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		err = c.roles.Create(ctx, role)
		if err != nil {
			return fmt.Errorf("failed creating role: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("failed getting role: %w", err)
	}
}
