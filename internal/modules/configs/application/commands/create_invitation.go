package commands

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
)

type CreateInvitationCommand struct {
	invitations   models.InvitationRepository
	organizations models.OrganizationRepository
	configs       models.ConfigRepository
	roles         models.RoleRepository
	groups        models.GroupRepository
	mail          *gomail.Dialer
}

func NewCreateInvitationCommand(
	invitations models.InvitationRepository,
	organizations models.OrganizationRepository,
	configs models.ConfigRepository,
	roles models.RoleRepository,
	groups models.GroupRepository,
	mail *gomail.Dialer,
) *CreateInvitationCommand {
	return &CreateInvitationCommand{
		invitations,
		organizations,
		configs,
		roles,
		groups,
		mail,
	}
}

type CreateInvitationRequest struct {
	ExpiresAt       time.Time
	OrganizationID  primitive.ObjectID
	CreatorID       primitive.ObjectID
	RoleGroupID     *primitive.ObjectID
	RolePermissions []string
	RoleName        string
	UserDisplayName string
	UserEmail       string
}

func (c *CreateInvitationCommand) CreateInvitation(ctx context.Context, req *CreateInvitationRequest) (*models.Invitation, error) {

	organization, err := c.organizations.GetByID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	config, err := c.configs.GetByOrganizationID(ctx, organization.ID)
	if err != nil {
		return nil, fmt.Errorf("config not found for organization: %w", err)
	}

	// validate role group if provided
	if req.RoleGroupID != nil {
		_, err = c.groups.GetByIDAndOrganizationID(ctx, *req.RoleGroupID, req.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("group not found in organization: %w", err)
		}
	}

	// validate creator permissions if creator is provided
	if !req.CreatorID.IsZero() {
		role, err := c.roles.GetByUserIDAndOrganizationID(ctx, req.CreatorID, req.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("creator role not found in organization: %w", err)
		}
		if !role.Permissions.Has("invitations.create") {
			return nil, fmt.Errorf("creator does not have permission to create invitations")
		}
	}

	invitation := &models.Invitation{
		OrganizationID:  req.OrganizationID,
		CreatorID:       req.CreatorID,
		ExpiresAt:       req.ExpiresAt,
		RoleGroupID:     req.RoleGroupID,
		RolePermissions: req.RolePermissions,
		RoleName:        req.RoleName,
		UserDisplayName: req.UserDisplayName,
		UserEmail:       req.UserEmail,
		CreatedAt:       time.Now(),
	}

	err = c.invitations.Create(ctx, invitation)
	if err != nil {
		return nil, fmt.Errorf("failed creating inviation: %w", err)
	}

	email, ok := config.Emails["invitation"]
	if !ok {
		return nil, fmt.Errorf("invitation email template not found in config")
	}

	data := struct {
		UserDisplayName  string
		OrganizationName string
		InvitationURL    string
		ConfigURL        string
		FileURL          models.ExternalFileURLFunc
	}{
		UserDisplayName:  req.UserDisplayName,
		OrganizationName: organization.Name,
		InvitationURL:    fmt.Sprintf("%s/invitation/%s", config.Url, invitation.ID.Hex()),
		ConfigURL:        config.Url,
		FileURL:          models.NewExternalFileURLFunc(config.Url, organization.ID),
	}

	subjectTemplate, err := template.New("subject").Parse(email.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed creating email subject: %w", err)
	}

	var subject bytes.Buffer
	err = subjectTemplate.Execute(&subject, data)
	if err != nil {
		return nil, fmt.Errorf("failed executing email subject template: %w", err)
	}

	bodyTemplate, err := template.New("body").Parse(email.Body)
	if err != nil {
		return nil, fmt.Errorf("failed creating email body: %w", err)
	}

	var body bytes.Buffer
	err = bodyTemplate.Execute(&body, data)
	if err != nil {
		return nil, fmt.Errorf("failed executing email body template: %w", err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", fmt.Sprintf("%s <%s>", organization.Name, config.Email))
	message.SetHeader("To", fmt.Sprintf("%s <%s>", req.UserDisplayName, req.UserEmail))
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", body.String())

	err = c.mail.DialAndSend(message)
	if err != nil {
		return nil, err
	}

	return invitation, nil

}
