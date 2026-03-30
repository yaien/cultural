package pages

import (
	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/web/dashboard/views/icons"
)

type SectionDefinition struct {
	Name      string
	Component func(s *State) templ.Component
	Icon      templ.Component
	Web       bool
	Delete    bool
	Page      bool
	Layout    bool
	Email     bool
	Tab       bool
}

const (
	InitialSection       = "initial"
	CreateSection        = "create"
	DeleteSection        = "delete"
	FilesSection         = "files"
	FileSection          = "file"
	FontsSection         = "fonts"
	BrowseFontsSection   = "browse-fonts"
	ConfigureFontSection = "configure-font"
	ColorsSection        = "colors"
	EditStylesSection    = "edit-stytles"
	EditScriptSection    = "edit-scripts"
	EditHTMLSection      = "edit-html"
	PublishSection       = "publish"
)

var Sections = []SectionDefinition{
	{Name: InitialSection, Component: Initial, Tab: true, Icon: icons.Home()},
	{Name: CreateSection, Component: Create, Tab: true, Web: true, Icon: icons.Add()},
	{Name: DeleteSection, Component: Delete, Tab: true, Delete: true, Icon: icons.Trash()},
	{Name: FilesSection, Component: Files, Tab: true, Icon: icons.Files()},
	{Name: FileSection, Component: File},
	{Name: FontsSection, Component: Fonts, Tab: true, Web: true, Icon: icons.BrandFamily()},
	{Name: BrowseFontsSection, Component: BrowseFonts},
	{Name: ConfigureFontSection, Component: ConfigureFont},
	{Name: ColorsSection, Component: Colors, Tab: true, Web: true, Icon: icons.Palette()},
	{Name: EditHTMLSection, Component: EditHTML, Tab: true, Icon: icons.Html()},
	{Name: EditStylesSection, Component: EditStyles, Tab: true, Web: true, Icon: icons.Css()},
	{Name: EditScriptSection, Component: EditScript, Tab: true, Web: true, Icon: icons.Js()},
	{Name: PublishSection, Component: Publish, Tab: true, Icon: icons.Publish()},
}

func Section(state *State) templ.Component {
	for _, section := range Sections {
		if section.Name == state.Section {
			return section.Component(state)
		}
	}

	return Initial(state)
}
