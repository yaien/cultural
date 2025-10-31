package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"github.com/yaien/cultural/internal/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
)

type CreateInvitationCommand struct {
	invitations models.InvitationRepository
	mail        *gomail.Dialer
}

func NewCreateInvitationCommand(invitations models.InvitationRepository, mail *gomail.Dialer) *CreateInvitationCommand {
	return &CreateInvitationCommand{
		invitations: invitations,
		mail:        mail,
	}
}

type CreateInvitationRequest struct {
	OrganizationID any
	GroupID        any
	CreatorID      any
	Email          string
	BaseURL        string
	Permissions    []string
	Name           string
	ExpiresAt      time.Time
}

func (c *CreateInvitationCommand) CreateInvitation(ctx context.Context, req *CreateInvitationRequest) (*models.Invitation, error) {
	invitation := &models.Invitation{
		ID:             primitive.NewObjectID(),
		OrganizationID: req.OrganizationID,
		GroupID:        req.GroupID,
		CreatorID:      req.CreatorID,
		Email:          req.Email,
		Permissions:    req.Permissions,
		Name:           req.Name,
		ExpiresAt:      req.ExpiresAt,
	}

	err := c.invitations.Create(ctx, invitation)
	if err != nil {
		return nil, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "example@example.com")
	message.SetHeader("To", invitation.Email)
	message.SetHeader("Subject", "You're invited!")
	message.SetBody("text/plain", fmt.Sprintf("You have been invited. Please use the following link %s/invited/%s to accept the invitation.", req.BaseURL, shared.IDToStr(invitation.ID)))

	err = c.mail.DialAndSend(message)
	if err != nil {
		return nil, err
	}

	return invitation, nil

}
