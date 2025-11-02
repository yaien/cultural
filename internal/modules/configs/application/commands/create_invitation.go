package commands

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/models"
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
	OrganizationID primitive.ObjectID
	GroupID        primitive.ObjectID
	CreatorID      primitive.ObjectID
	Email          string
	BaseURL        string
	Permissions    []string
	Name           string
	DisplayName    string
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
		DisplayName:    req.DisplayName,
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

	err = c.mail.DialAndSend(message)
	if err != nil {
		return nil, err
	}

	return invitation, nil

}
