package models

import (
	"fmt"
	"html/template"
)

type Page struct {
	Title  string        `bson:"title,omitempty" json:"title,omitempty"`
	Name   string        `bson:"name,omitempty" json:"name,omitempty"`
	Styles template.CSS  `bson:"styles,omitempty" json:"styles,omitempty"`
	Body   template.HTML `bson:"body" json:"body"`
}

var PageTemplate = template.Must(template.New("page").Parse(read("templates/page.html")))

type pageDataOptions struct {
	InlineStyles bool
	FilePath     string
	Page         *Page
	Fonts        map[string]*Font
	Colors       map[string]string
	Version      int64
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

func (p *pageDataOptions) WithVersion(version int64) *pageDataOptions {
	p.Version = version
	return p
}

type pageData struct {
	InlineStyles bool
	FilePath     string
	Page         *Page
	Fonts        map[string]*Font
	Colors       map[string]string
	Version      int64
}

func (p *pageDataOptions) Data() *pageData {
	return &pageData{
		InlineStyles: p.InlineStyles,
		FilePath:     p.FilePath,
		Page:         p.Page,
		Version:      p.Version,
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
