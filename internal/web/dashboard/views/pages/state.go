package pages

import (
	"fmt"

	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/application/storage"
)

type SelectedType = label.DraftModelType

const (
	SelectedTypePage   = label.DraftPageModelType
	SelectedTypeLayout = label.DraftLayoutModelType
	SelectedTypeEmail  = label.DraftEmailModelType
)

const (
	DefaultPageName   = label.DefaultPageName
	DefaultLayoutName = label.DefaultLayoutName
	DefaultEmailName  = label.DefaultEmailName
)

type State struct {
	Config             *label.Config
	Draft              *label.Draft
	Selected           any
	SelectedType       SelectedType
	SelectedKey        string
	SelectedFileName   string
	SelectedFontFamily string
	SelectedFontKey    string
	FileURL            FileURLFunc
	File               FileFunc
	Files              FilesFunc
	Font               FontFunc
	Fonts              FontsFunc
	Section            string
}

type FileFunc func(filename string) (storage.File, error)
type FontFunc func(fontFamily string) (label.Font, error)
type FontsFunc func(family string, limit, offset int) ([]label.Font, error)
type FilesFunc func() ([]storage.File, error)
type FileURLFunc storage.URLFunc

func (c *State) PageIsDefault() bool {
	page, ok := c.Selected.(*label.Page)
	if !ok {
		return false
	}

	return page.Name == DefaultPageName
}

func (c *State) PageUrl() string {
	page, ok := c.Selected.(*label.Page)
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

func (c *State) HxSelectedVals() string {
	return fmt.Sprintf("{%q: %q, %q: %q}", SelectedTypeQuery, c.SelectedType, SelectedKeyQuery, c.SelectedKey)
}

func (c *State) HxSelectedSectionVals() string {
	return fmt.Sprintf("{%q: %q, %q: %q, %q: %q}", SelectedTypeQuery, c.SelectedType, SelectedKeyQuery, c.SelectedKey, SectionQuery, c.Section)
}
