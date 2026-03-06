package pages

import (
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type SelectedType string

const (
	SelectedTypePage   SelectedType = "page"
	SelectedTypeLayout SelectedType = "layout"
	SelectedTypeEmail  SelectedType = "email"
)

const (
	DefaultPageName   = "index"
	DefaultLayoutName = "default"
	DefaultEmailName  = "invitation"
)

type State struct {
	Config             *models.Config
	Draft              *models.Draft
	Selected           any
	SelectedType       SelectedType
	SelectedKey        string
	SelectedFileName   string
	SelectedFontFamily string
	FileURL            FileURLFunc
	File               FileFunc
	Files              FilesFunc
	Font               FontFunc
	Section            string
}

type FileFunc func(filename string) (*models.File, error)
type FontFunc func(fontFamily string) (*models.Font, error)
type FilesFunc func() ([]*models.File, error)
type FileURLFunc models.FileURLFunc

func (c *State) PageIsDefault() bool {
	page, ok := c.Selected.(*models.Page)
	if !ok {
		return false
	}

	return page.Name == DefaultPageName
}

func (c *State) PageUrl() string {
	page, ok := c.Selected.(*models.Page)
	if !ok {
		return ""
	}

	if page.Name == DefaultPageName {
		return c.Config.Url
	}

	return c.Config.Url + "/" + page.Name
}

func (c *State) NotDeleteable() bool {
	switch c.SelectedType {
	case SelectedTypePage:
		return c.SelectedKey == DefaultPageName
	case SelectedTypeLayout:
		return c.SelectedKey == DefaultLayoutName
	default:
		return true
	}

}

func (c *State) NotWeb() bool {
	return c.SelectedType == SelectedTypeEmail
}

func (c *State) NotPage() bool {
	return c.SelectedType != SelectedTypePage
}
