package models

import (
	"embed"
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

const (
	DefaultPageName   = "index"
	DefaultLayoutName = "default"
	DefaultEmailName  = "invitation"
)

var DefaultLayout = &Layout{
	Name:  DefaultLayoutName,
	Title: "Diseño por defecto",
	Body:  `{{template "page_body" .}}`,
}

var DefaultLayouts = map[string]*Layout{
	DefaultLayoutName: DefaultLayout,
}

var DefaultPages = map[string]*Page{
	DefaultPageName: {
		Name:   DefaultPageName,
		Styles: read("templates/index_page.css"),
		Body:   read("templates/index_page.html"),
	},
}

var DefaultEmails = map[string]*Email{
	DefaultEmailName: {
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
