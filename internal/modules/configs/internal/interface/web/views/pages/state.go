package pages

import (
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type State struct {
	Config             *models.Config
	Draft              *models.Draft
	Selected           any
	SelectedType       string
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

func (c *State) PageIsIndex() bool {
	page, ok := c.Selected.(*models.Page)
	if !ok {
		return false
	}

	return page.Name == "index"
}

func (c *State) PageUrl() string {
	page, ok := c.Selected.(*models.Page)
	if !ok {
		return ""
	}

	if page.Name == "index" {
		return c.Config.Url
	}

	return c.Config.Url + "/" + page.Name
}

func (c *State) SelectedIsDeleteable() bool {
	switch sel := c.Selected.(type) {
	case *models.Page:
		return sel.Name != "index"
	case *models.Layout:
		return sel.Name != "default"
	default:
		return false
	}
}

func (c *State) SelectedIsForWeb() bool {
	switch c.Selected.(type) {
	case *models.Page, *models.Layout:
		return true
	default:
		return false
	}
}
