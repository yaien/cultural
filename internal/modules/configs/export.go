package configs

import (
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type CreateInvitationRequest = commands.CreateInvitationRequest
type GetFileRequest = queries.GetFileRequest
type GetFileResponse = queries.GetFileResponse

type Config = models.Config

const ConfigContextKey = models.ConfigContextKey

var NewPageData = models.NewPageData
var WritePageBaseStyles = models.WritePageBaseStyles
var PageTemplate = models.PageTemplate
var EmptyPage = models.EmptyPage
var DefaultLayout = models.DefaultLayout
