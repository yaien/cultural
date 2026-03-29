package configs

import (
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type CreateInvitationRequest = commands.CreateInvitationRequest

type Config = models.Config
type PageData = models.PageData

const ConfigContextKey = middlewares.ConfigContextKey

var WritePageBaseStyles = models.WritePageBaseStyles
var EmptyPage = models.EmptyPage
var DefaultLayout = models.DefaultLayout
var RenderPage = models.RenderPage
