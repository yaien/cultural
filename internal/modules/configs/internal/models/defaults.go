package models

import (
	"embed"
	"html/template"
)

//go:embed templates/*
var fs embed.FS

func read(file string) string {
	content, err := fs.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(content)
}

var DefaultLayout = &Page{
	Name:  "default",
	Title: "Diseño por defecto",
	Body:  `{{template "page_body" .}}`,
}

var DefaultLayouts = map[string]*Page{
	"default": DefaultLayout,
}

var DefaultPages = map[string]*Page{
	"index": {
		Title:  "Inicio",
		Name:   "index",
		Styles: template.CSS(read("templates/index_page.css")),
		Body:   template.HTML(read("templates/index_page.html")),
	},
}

var DefaultEmails = map[string]*Email{
	"invitation": {
		Subject: read("templates/invitation_email_subject.txt"),
		Body:    read("templates/invitation_email_body.html"),
	},
}

var DefaultColors = map[string]string{
	"primary":    "#330136",
	"secondary":  "#FFFFFF",
	"accent":     "#FF6F61",
	"background": "#F5F5F5",
	"text":       "#333333",
}
