package models

import (
	"bytes"
	"fmt"
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

func (c *pageData) Styles() (template.HTML, error) {
	buff := &bytes.Buffer{}
	data := &PageStyleTemplateData{
		Fonts:  c.Fonts,
		Colors: c.Colors,
		Styles: c.Page.Styles,
	}

	if err := PageStyleTemplate.Execute(buff, data); err != nil {
		return "", err
	}
	styles := fmt.Sprintf("<style type=%q>\n%s</style>", "text/css", buff.String())

	return template.HTML(styles), nil
}
