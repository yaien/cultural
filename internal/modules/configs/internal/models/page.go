package models

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type Layout struct {
	Title  string `bson:"title,omitempty" json:"title,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Layout string `bson:"layout,omitempty" json:"layout,omitempty"`
	Styles string `bson:"styles,omitempty" json:"styles,omitempty"`
	Script string `bson:"script,omitempty" json:"script,omitempty"`
	Body   string `bson:"body" json:"body"`
}

type Page struct {
	Title       string `bson:"title,omitempty" json:"title,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Name        string `bson:"name,omitempty" json:"name,omitempty"`
	Layout      string `bson:"layout,omitempty" json:"layout,omitempty"`
	Styles      string `bson:"styles,omitempty" json:"styles,omitempty"`
	Script      string `bson:"script,omitempty" json:"script,omitempty"`
	Body        string `bson:"body" json:"body"`
	OGImage     string `bson:"ogImage,omitempty" json:"ogImage,omitempty"`
	OGType      string `bson:"ogType,omitempty" json:"ogType,omitempty"`
}

var EmptyPage = &Page{}

type PageData struct {
	InlineStyles        bool
	InlineScript        bool
	FileURLFunc         FileURLFunc
	ExternalFileURLFunc FileURLFunc
	AppTitle            string
	Page                *Page
	Layout              *Layout
	Fonts               Fonts
	Colors              Colors
	Version             int64
}

func (c *PageData) FileURL(name string, variant ...int) string {
	return c.FileURLFunc(name, variant...)
}

func (c *PageData) ExternalFileURL(name string, variant ...int) string {
	return c.ExternalFileURLFunc(name, variant...)
}

// Title returns the full title of the page, combining the page title and app title.
func (p *PageData) Title() string {
	if p.Page.Title != "" {
		return p.Page.Title + " - " + p.AppTitle
	}

	return p.AppTitle
}

type pageStyleTemplateData struct {
	Fonts        map[string]*Font
	Colors       Colors
	PageStyles   string
	LayoutStyles string
}

var pageStyleTemplate = template.Must(template.New("styles").Parse(read("templates/styles.txt")))

// WritePageBaseStyles writes the base styles for a page to the provided writer.
func WritePageBaseStyles(b io.Writer, cfg *Config) error {
	return pageStyleTemplate.Execute(b, &pageStyleTemplateData{
		Fonts:  cfg.Fonts,
		Colors: cfg.Colors,
	})
}

// Styles generates the combined styles for the page and layout, including fonts and colors.
func (c *PageData) Styles() (string, error) {
	buff := &bytes.Buffer{}
	data := &pageStyleTemplateData{
		Fonts:        c.Fonts,
		Colors:       c.Colors,
		PageStyles:   c.Page.Styles,
		LayoutStyles: c.Layout.Styles,
	}

	if err := pageStyleTemplate.Execute(buff, data); err != nil {
		return "", err
	}
	styles := fmt.Sprintf("<style type=%q>\n%s</style>", "text/css", buff.String())

	return styles, nil
}

var pageScriptTemplate = template.Must(template.New("script").Parse(read("templates/scripts.txt")))

type pageScriptTemplateData struct {
	PageScript   string
	LayoutScript string
}

// Script generates the combined scripts for the page and layout.
func (c *PageData) Script() (string, error) {
	data := &pageScriptTemplateData{
		PageScript:   string(c.Page.Script),
		LayoutScript: string(c.Layout.Script),
	}

	var buff bytes.Buffer
	if err := pageScriptTemplate.Execute(&buff, data); err != nil {
		return "", fmt.Errorf("failed executing script template: %w", err)
	}

	return buff.String(), nil
}

// FontLinks generates the necessary <link> tags for the fonts used in the page.
func (p *PageData) FontLinks() (string, error) {
	setted := make(map[string]bool)
	sb := strings.Builder{}
	for _, font := range p.Fonts {
		switch font.Provider {
		case "google":
			if !setted[font.Provider] {
				setted[font.Provider] = true
				fmt.Fprintln(&sb, `<link rel="preconnect" href="https://fonts.googleapis.com"/>`)
				fmt.Fprintln(&sb, `<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>`)
			}
			fmt.Fprintf(&sb, `<link rel="stylesheet" href="%s"/>`, googleFontURL(font))
		}
	}
	return sb.String(), nil
}

func googleFontURL(font *Font) string {
	url := "https://fonts.googleapis.com/css2?"
	family := strings.ReplaceAll(font.Family, " ", "+")
	url += "family=" + family

	// Manejar variantes correctamente
	if len(font.Variants) > 0 {
		weights := []string{}
		italicWeights := []string{}
		hasItalic := false

		for _, variant := range font.Variants {
			switch {
			case variant == "italic":
				hasItalic = true
				italicWeights = append(italicWeights, "400")
			case strings.HasSuffix(variant, "italic"):
				hasItalic = true
				weight := strings.TrimSuffix(variant, "italic")
				if weight == "" {
					weight = "400"
				}
				italicWeights = append(italicWeights, weight)
			case variant == "regular":
				weights = append(weights, "400")
			default:
				// Peso numérico (300, 400, 500, etc.)
				weights = append(weights, variant)
			}
		}

		// Si no hay pesos específicos, agregar 400 por defecto
		if len(weights) == 0 && !hasItalic {
			weights = append(weights, "400")
		}

		// Construir URL según las variantes
		if hasItalic && len(weights) > 0 {
			// Formato: ital,wght@0,400;0,700;1,400;1,700
			url += ":ital,wght@"
			var params []string

			// Pesos normales
			for _, w := range weights {
				params = append(params, "0,"+w)
			}

			// Pesos itálicos (usar weights si italicWeights está vacío)
			targetItalicWeights := italicWeights
			if len(targetItalicWeights) == 0 && len(weights) > 0 {
				targetItalicWeights = weights
			}

			for _, w := range targetItalicWeights {
				params = append(params, "1,"+w)
			}

			url += strings.Join(params, ";")
		} else if hasItalic {
			// Solo itálica sin pesos específicos
			url += ":ital,wght@1,400"
		} else if len(weights) > 0 {
			// Solo pesos sin itálica
			url += ":wght@" + strings.Join(weights, ";")
		}
	} else {
		// Sin variantes especificadas - cargar básicas
		url += ":wght@300;400;500;600;700"
	}

	// Agregar display=swap para mejor rendimiento
	url += "&display=swap"
	return url
}

var pageTemplate = template.Must(template.New("page").Parse(read("templates/page.html")))

// This file contains the logic for rendering a page using the page and layout templates.
func RenderPage(data *PageData) (string, error) {
	var buffer bytes.Buffer

	base, err := pageTemplate.Clone()
	if err != nil {
		return "", fmt.Errorf("failed decoding template: %w", err)
	}

	if data.Page == nil {
		return "", fmt.Errorf("page data is nil")
	}

	if data.Layout == nil {
		return "", fmt.Errorf("layout data is nil")
	}

	parsed, err := base.Parse(fmt.Sprintf(`{{define "layout_body"}}%s{{end}}{{define "page_body"}}%s{{end}}`, data.Layout.Body, data.Page.Body))
	if err != nil {
		return "", fmt.Errorf("failed parsing template: %w", err)
	}

	if err := parsed.Execute(&buffer, data); err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buffer.String(), nil
}
