package models

import (
	"fmt"
	"html/template"
)

type Page struct {
	Title  string        `bson:"title,omitempty" json:"title,omitempty"`
	Name   string        `bson:"name,omitempty" json:"name,omitempty"`
	Layout string        `bson:"layout,omitempty" json:"layout,omitempty"`
	Styles template.CSS  `bson:"styles,omitempty" json:"styles,omitempty"`
	Script template.JS   `bson:"script,omitempty" json:"script,omitempty"`
	Body   template.HTML `bson:"body" json:"body"`
}

var PageTemplate = template.Must(template.New("page").Parse(read("templates/page.html")))

var EmptyPage = &Page{}

type pageDataOptions struct {
	InlineStyles bool
	InlineScript bool
	FilePath     string
	Page         *Page
	Layout       *Page
	Fonts        map[string]*Font
	Colors       map[string]string
	AppTitle     string
	Version      int64
}

func NewPageData(page *Page, layout *Page) *pageDataOptions {
	return &pageDataOptions{Page: page, Layout: layout}
}

func (p *pageDataOptions) WithAppTitle(appTitle string) *pageDataOptions {
	p.AppTitle = appTitle
	return p
}

func (p *pageDataOptions) WithInlineStyles(inlineStyles bool) *pageDataOptions {
	p.InlineStyles = inlineStyles
	return p
}

func (p *pageDataOptions) WithInlineScript(inlineScript bool) *pageDataOptions {
	p.InlineScript = inlineScript
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

func (p *pageDataOptions) WithVersion(version int64) *pageDataOptions {
	p.Version = version
	return p
}

type pageData struct {
	InlineStyles bool
	InlineScript bool
	FilePath     string
	AppTitle     string
	Page         *Page
	Layout       *Page
	Fonts        map[string]*Font
	Colors       map[string]string
	Version      int64
}

func (p *pageDataOptions) Data() *pageData {
	return &pageData{
		InlineStyles: p.InlineStyles,
		InlineScript: p.InlineScript,
		FilePath:     p.FilePath,
		Page:         p.Page,
		Layout:       p.Layout,
		Version:      p.Version,
		AppTitle:     p.AppTitle,
		Fonts:        p.Fonts,
		Colors:       p.Colors,
	}
}

func (p *pageData) FileURL(filename string, variant ...int) string {
	if len(variant) > 0 {
		return fmt.Sprintf("%s/%s?variant=%d", p.FilePath, filename, variant[0])
	}
	return p.FilePath + "/" + filename
}

func (p *pageData) Title() string {
	if p.Page.Title != "" {
		return p.Page.Title + " - " + p.AppTitle
	}

	return p.AppTitle
}
