package configs

import (
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type CreateInvitationRequest = commands.CreateInvitationRequest
type Config = models.Config

const ConfigContextKey = models.ConfigContextKey

var NewPageData = models.NewPageData
var WritePageBaseStyles = models.WritePageBaseStyles
var PageTemplate = models.PageTemplate
