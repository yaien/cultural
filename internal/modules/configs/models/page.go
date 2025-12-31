package models

import (
	"html/template"
)

type Page struct {
	Title  string        `bson:"title,omitempty" json:"title,omitempty"`
	Name   string        `bson:"name,omitempty" json:"name,omitempty"`
	Styles template.CSS  `bson:"styles,omitempty" json:"styles,omitempty"`
	Body   template.HTML `bson:"body" json:"body"`
}

var PageBaseStyles = template.Must(template.New("styles").Parse(`
	:root {
	{{range $key, $font := .Fonts}}
		--font-{{ $key }}: '{{ $font.Family }}', sans-serif;
	{{ end }}
	{{range $key, $value := .Colors}}
		--color-{{ $key }}: {{ $value }};
	{{ end }}
	}
`))

var PageTemplate = template.Must(template.New("page").Parse(read("templates/page.html")))

type pageData struct {
	InlineStyles bool
	FilePath     string
	Page         *Page
	Config       *Config
	Components   *pageComponents
}

type pageDataOptions struct {
	InlineStyles bool
	FilePath     string
	Page         *Page
	Config       *Config
}

func NewPageData(config *Config, page *Page) *pageDataOptions {
	return &pageDataOptions{Config: config, Page: page}
}

func (p *pageDataOptions) WithInlineStyles(inlineStyles bool) *pageDataOptions {
	p.InlineStyles = inlineStyles
	return p
}

func (p *pageDataOptions) WithFilePath(filepath string) *pageDataOptions {
	p.FilePath = filepath
	return p
}

func (p *pageDataOptions) Data() *pageData {
	return &pageData{
		InlineStyles: p.InlineStyles,
		FilePath:     p.FilePath,
		Page:         p.Page,
		Config:       p.Config,
		Components:   &pageComponents{options: p},
	}
}

type pageComponents struct {
	options *pageDataOptions
}
