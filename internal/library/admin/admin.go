package admin

import "github.com/yaien/cultural/internal/library/mail"

type Admin struct {
	Invitations *Invitations
	Roles       *Roles
}

func New(roles RoleRepository, organizations OrganizationRepository, invitations InvitationRepository, groups GroupRepository, m mail.Mail) *Admin {
	return &Admin{
		Invitations: NewInvitations(roles, organizations, groups, invitations, m),
		Roles:       NewRoles(roles),
	}
}
