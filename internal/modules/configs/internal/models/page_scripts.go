package models

import (
	"bytes"
	"fmt"
	"html/template"
)

var PageScriptTemplate = template.Must(template.New("script").Parse(read("templates/scripts.txt")))

type PageScriptTemplateData struct {
	PageScript   template.JS
	LayoutScript template.JS
}

func (c *pageData) Script() (template.HTML, error) {
	data := &PageScriptTemplateData{
		PageScript:   c.Page.Script,
		LayoutScript: c.Layout.Script,
	}

	var buff bytes.Buffer
	if err := PageScriptTemplate.Execute(&buff, data); err != nil {
		return "", fmt.Errorf("failed executing script template: %w", err)
	}

	return template.HTML(buff.String()), nil
}
