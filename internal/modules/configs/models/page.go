package models

import (
	"html/template"

	"github.com/a-h/templ"
)

type Node struct {
	Type     string           `bson:"type,omitempty" json:"type,omitempty"`
	Attrs    templ.Attributes `bson:"attrs,omitempty" json:"attrs,omitempty"`
	Content  string           `bson:"content,omitempty" json:"content,omitempty"`
	Children []Node           `bson:"children,omitempty" json:"children,omitempty"`
}

type Page struct {
	Title  string `bson:"title,omitempty" json:"title,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Styles string `bson:"styles,omitempty" json:"styles,omitempty"`
	Body   Node   `bson:"body" json:"body"`
}

var Styles = template.Must(template.New("styles").Parse(`
	:root {
	{{range $key, $font := .Fonts}}
		--font-{{ $key }}: '{{ $font.Family }}', sans-serif;
	{{ end }}		
	{{range $key, $value := .Colors}}
		--color-{{ $key }}: {{ $value }};
	{{ end }}
	}
`))
