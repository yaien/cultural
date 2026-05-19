package pages

import (
	"github.com/a-h/templ"
)

type SectionDefinition struct {
	Name      string
	Component func(s *State) templ.Component
	Icon      string
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
	{Name: InitialSection, Component: Initial, Tab: true, Icon: "fa-solid fa-house"},
	{Name: CreateSection, Component: Create, Tab: true, Web: true, Icon: "fa-solid fa-plus"},
	{Name: DeleteSection, Component: Delete, Tab: true, Delete: true, Icon: "fa-solid fa-trash-can"},
	{Name: FilesSection, Component: Files, Tab: true, Icon: "fa-solid fa-image"},
	{Name: FileSection, Component: File},
	{Name: FontsSection, Component: Fonts, Tab: true, Web: true, Icon: "fa-solid fa-font"},
	{Name: BrowseFontsSection, Component: BrowseFonts},
	{Name: ConfigureFontSection, Component: ConfigureFont},
	{Name: ColorsSection, Component: Colors, Tab: true, Web: true, Icon: "fa-solid fa-palette"},
	{Name: EditHTMLSection, Component: EditHTML, Tab: true, Icon: "fa-solid fa-code"},
	{Name: EditStylesSection, Component: EditStyles, Tab: true, Web: true, Icon: "fa-brands fa-css"},
	{Name: EditScriptSection, Component: EditScript, Tab: true, Web: true, Icon: "fa-brands fa-js"},
	{Name: PublishSection, Component: Publish, Tab: true, Icon: "fa-solid fa-upload"},
}

func Section(state *State) templ.Component {
	for _, section := range Sections {
		if section.Name == state.Section {
			return section.Component(state)
		}
	}

	return Initial(state)
}
