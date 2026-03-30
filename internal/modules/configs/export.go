package configs

import (
	"github.com/yaien/cultural/internal/label"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
)

type Config = label.Config
type PageData = label.PageData

const ConfigContextKey = middlewares.ConfigContextKey

var WritePageBaseStyles = label.WritePageBaseStyles
var EmptyPage = label.EmptyPage
var DefaultLayout = label.DefaultLayout
var RenderPage = label.RenderPage
