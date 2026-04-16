package admin

import (
	"github.com/yaien/cultural/internal/lib/mail"
	"gorm.io/gorm"
)

type Admin struct {
	Invitations *Invitations
	Roles       *Roles
}

func New(db *gorm.DB, m mail.Mail) *Admin {
	return &Admin{
		Invitations: NewInvitations(db, m),
		Roles:       NewRoles(db),
	}
}
