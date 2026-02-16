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

var PageTemplate = template.Must(template.New("page").Parse(read("templates/page.html")))

type pageData struct {
	InlineStyles bool
	FilePath     string
	Page         *Page
	Fonts        map[string]*Font
	Colors       map[string]string
	Components   *pageComponents
}

type pageDataOptions struct {
	InlineStyles bool
	FilePath     string
	Page         *Page
	Fonts        map[string]*Font
	Colors       map[string]string
}

func NewPageData(page *Page) *pageDataOptions {
	return &pageDataOptions{Page: page}
}

func (p *pageDataOptions) WithInlineStyles(inlineStyles bool) *pageDataOptions {
	p.InlineStyles = inlineStyles
	return p
}

func (p *pageDataOptions) WithFilePath(filepath string) *pageDataOptions {
	p.FilePath = filepath
	return p
}

func (p *pageDataOptions) WithFonts(fonts map[string]*Font) *pageDataOptions {
	p.Fonts = fonts
	return p
}

func (p *pageDataOptions) WithColors(colors map[string]string) *pageDataOptions {
	p.Colors = colors
	return p
}

func (p *pageDataOptions) Data() *pageData {
	return &pageData{
		InlineStyles: p.InlineStyles,
		FilePath:     p.FilePath,
		Page:         p.Page,
		Components:   &pageComponents{options: p},
	}
}

type pageComponents struct {
	options *pageDataOptions
}
