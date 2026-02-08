package models

import (
	"fmt"
	"html/template"
	"io"
)

type Page struct {
	Title  string        `bson:"title,omitempty" json:"title,omitempty"`
	Name   string        `bson:"name,omitempty" json:"name,omitempty"`
	Styles template.CSS  `bson:"styles,omitempty" json:"styles,omitempty"`
	Body   template.HTML `bson:"body" json:"body"`
}

func WritePageBaseStyles(b io.Writer, cfg *Config) error {

	_, err := fmt.Fprintln(b, ":root {")
	if err != nil {
		return err
	}

	for key, font := range cfg.Fonts {
		_, err := fmt.Fprintf(b, "\t--font-%s: %q, sans-serif;\n", key, font.Family)
		if err != nil {
			return err
		}
	}
	for key, value := range cfg.Colors {
		_, err := fmt.Fprintf(b, "\t--color-%s: %s;\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(b, "}")
	if err != nil {
		return err
	}

	return nil
}

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
