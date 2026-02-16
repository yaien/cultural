package models

import (
	"bytes"
	"html/template"
	"io"
)

var PageStyleTemplate = template.Must(template.New("styles").Parse(read("templates/styles.txt")))

type PageStyleTemplateData struct {
	Fonts  map[string]*Font
	Colors map[string]string
	Styles template.CSS
}

func WritePageBaseStyles(b io.Writer, cfg *Config) error {
	return PageStyleTemplate.Execute(b, &PageStyleTemplateData{
		Fonts:  cfg.Fonts,
		Colors: cfg.Colors,
	})
}

func (c *pageComponents) Styles() (template.CSS, error) {
	buff := &bytes.Buffer{}
	data := &PageStyleTemplateData{
		Fonts:  c.options.Fonts,
		Colors: c.options.Colors,
		Styles: c.options.Page.Styles,
	}

	if err := PageStyleTemplate.Execute(buff, data); err != nil {
		return "", err
	}

	return template.CSS(buff.String()), nil
}
