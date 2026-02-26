package models

import (
	"bytes"
	"fmt"
	"text/template"
)

var PageScriptTemplate = template.Must(template.New("script").Parse(read("templates/scripts.txt")))

type PageScriptTemplateData struct {
	PageScript   string
	LayoutScript string
}

func (c *pageData) Script() (string, error) {
	data := &PageScriptTemplateData{
		PageScript:   string(c.Page.Script),
		LayoutScript: string(c.Layout.Script),
	}

	var buff bytes.Buffer
	if err := PageScriptTemplate.Execute(&buff, data); err != nil {
		return "", fmt.Errorf("failed executing script template: %w", err)
	}

	return buff.String(), nil
}
