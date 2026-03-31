package admin

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"

	"github.com/yaien/cultural/internal/application/auth"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/mail"
)

type Invitation struct {
	ID              primitive.ID `gorm:"primaryKey;autoIncrement"`
	OrganizationID  primitive.ID `gorm:"index"`
	Organization    *Organization
	CreatorID       primitive.ID
	Creator         *auth.User
	CreatedAt       time.Time
	AcceptedAt      *time.Time
	ExpiresAt       time.Time
	RoleGroupID     *primitive.ID
	RoleGroup       *Group
	RolePermissions Permissions `gorm:"type:jsonb;serializer:json"`
	RoleName        string
	UserDisplayName string
	UserEmail       string `gorm:"index"`
}

type InvitationRepository interface {
	GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ID) (*Invitation, error)
	Create(ctx context.Context, invitation *Invitation) error
	Update(ctx context.Context, invitation *Invitation) error
}

type Invitations struct {
	roles         RoleRepository
	organizations OrganizationRepository
	groups        GroupRepository
	invitations   InvitationRepository
	mail          mail.Mail
}

func NewInvitations(roles RoleRepository, organizations OrganizationRepository, groups GroupRepository, invitations InvitationRepository, mail mail.Mail) *Invitations {
	return &Invitations{
		roles:         roles,
		organizations: organizations,
		groups:        groups,
		invitations:   invitations,
		mail:          mail,
	}
}

type CreateInvitationOptions struct {
	ExpiresAt       time.Time
	OrganizationID  primitive.ID
	CreatorID       primitive.ID
	RoleGroupID     *primitive.ID
	Config          *label.Config
	RolePermissions []string
	RoleName        string
	UserDisplayName string
	UserEmail       string
}

type InvitationEmailData struct {
	UserDisplayName  string
	OrganizationName string
	InvitationURL    string
	ConfigURL        string
	FileURL          storage.URLFunc
}

func (c *Invitations) Create(ctx context.Context, req *CreateInvitationOptions) (*Invitation, error) {

	organization, err := c.organizations.GetByID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	// validate role group if provided
	if req.RoleGroupID != nil {
		_, err = c.groups.GetByIDAndOrganizationID(ctx, *req.RoleGroupID, req.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("group not found in organization: %w", err)
		}
	}

	// validate creator permissions if creator is provided
	if req.CreatorID != 0 {
		role, err := c.roles.GetByUserIDAndOrganizationID(ctx, req.CreatorID, req.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("creator role not found in organization: %w", err)
		}
		if !role.Permissions.Has("invitations.create") {
			return nil, fmt.Errorf("creator does not have permission to create invitations")
		}
	}

	invitation := &Invitation{
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

	email, ok := req.Config.Emails["invitation"]
	if !ok {
		return nil, fmt.Errorf("invitation email template not found in config")
	}

	data := InvitationEmailData{
		UserDisplayName:  req.UserDisplayName,
		OrganizationName: organization.Name,
		InvitationURL:    fmt.Sprintf("%s/invitation/%d", req.Config.Url, invitation.ID),
		ConfigURL:        req.Config.Url,
		FileURL:          storage.NewExternalURLFunc(req.Config.Url, organization.ID),
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

	err = c.mail.Send(ctx, &mail.Email{
		To:       mail.Recipient{Name: req.UserDisplayName, Email: req.UserEmail},
		From:     mail.Recipient{Name: organization.Name, Email: req.Config.Email},
		Subject:  subject.String(),
		Body:     body.String(),
		Category: "invitation",
	})

	if err != nil {
		return nil, fmt.Errorf("failed sending the invitation email: %w", err)
	}

	return invitation, nil

}

type AcceptInvitationOptions struct {
	InvitationID   primitive.ID
	OrganizationID primitive.ID
	UserID         primitive.ID
	UserEmail      string
	UserName       string
	UserAvatarUrl  string
}

func (c *Invitations) Accept(ctx context.Context, req *AcceptInvitationOptions) error {
	invitation, err := c.invitations.GetByIDAndOrganizationID(ctx, req.InvitationID, req.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed getting invitation: %w", err)
	}

	if invitation.UserEmail != req.UserEmail {
		return coderror.Newf("email_missmatch", "%s user email doesn't match with this invitation", req.UserEmail)
	}

	if invitation.AcceptedAt != nil {
		return coderror.Newf("invitation_already_accepted", "invitation %d already accepted", req.InvitationID)
	}

	now := time.Now()
	invitation.AcceptedAt = &now

	err = c.invitations.Update(ctx, invitation)
	if err != nil {
		return fmt.Errorf("failed updating invitation: %w", err)
	}

	role, err := c.roles.GetByUserIDAndOrganizationID(ctx, req.UserID, req.OrganizationID)

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

	case coderror.Is(err, coderror.NotFound):
		role = &Role{
			UserID:         req.UserID,
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
