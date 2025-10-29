package models

import "github.com/a-h/templ"

type Node struct {
	Type     string            `bson:"type,omitempty" json:"type,omitempty"`
	Style    map[string]string `bson:"style,omitempty" json:"style,omitempty"`
	Attrs    templ.Attributes  `bson:"attrs,omitempty" json:"attrs,omitempty"`
	Content  string            `bson:"content,omitempty" json:"content,omitempty"`
	Children []*Node           `bson:"children,omitempty" json:"children,omitempty"`
}

type Fonts struct {
	Type     string            `bson:"type,omitempty" json:"type,omitempty"`
	Families map[string]string `bson:"families,omitempty" json:"families,omitempty"`
}

type Page struct {
	Name   string  `bson:"name,omitempty" json:"name,omitempty"`
	Styles string  `bson:"styles,omitempty" json:"styles,omitempty"`
	Body   []*Node `bson:"body,omitempty" json:"body,omitempty"`
}

var DefaultFonts = &Fonts{
	Type: "google",
	Families: map[string]string{
		"headings": "Instrument Serif",
		"primary":  "Outfit",
	},
}

var DefaultColors = map[string]string{
	"primary":    "#330136",
	"secondary":  "#FFFFFF",
	"accent":     "#FF6F61",
	"background": "#F5F5F5",
	"text":       "#333333",
}

var DefaultSite = &Page{
	Name: "index",
	Body: []*Node{
		{
			Type: "header",
			Children: []*Node{
				{
					Type: "nav",
					Children: []*Node{
						{Type: "a", Attrs: map[string]any{"href": "#festival"}, Content: "Festival"},
						{Type: "a", Attrs: map[string]any{"href": "#dates"}, Content: "Dates"},
						{Type: "a", Attrs: map[string]any{"href": "#tickets"}, Content: "Tickets"},
						{Type: "a", Attrs: map[string]any{"href": "#locations"}, Content: "Locations"},
						{Type: "a", Attrs: map[string]any{"href": "#about"}, Content: "About"},
					},
				},
			},
		},
		{
			Type: "main",
			Children: []*Node{
				{
					Type:  "section",
					Style: map[string]string{"text-align": "center"},
					Children: []*Node{
						{
							Type:  "img",
							Style: map[string]string{"max-width": "100%"},
							Attrs: map[string]any{"src": "/assets/static/landing/brand.png"},
						},
						{
							Type:    "h2",
							Content: "INTERNATIONAL DANCE, MUSIC & CULTURE FESTIVAL",
						},
					},
				},
				{
					Type:  "section",
					Attrs: map[string]any{"id": "festival"},
					Children: []*Node{
						{
							Type:    "h1",
							Content: "Lorem Ipsum",
						},
						{
							Type: "p",
							Content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
								Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
								Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. 
								Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. 
								Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
								`,
						},
						{
							Type:  "carousel",
							Style: map[string]string{"margin": "1rem auto"},
							Children: []*Node{
								{Type: "div", Content: "random pic 1"},
								{Type: "div", Content: "random pic 2"},
								{Type: "div", Content: "random pic 3"},
							},
						},
						{
							Type:  "section",
							Style: map[string]string{"text-align": "center"},
							Children: []*Node{
								{
									Type:    "raw",
									Content: `<iframe width="560" height="315" src="https://www.youtube.com/embed/Y4HWvsGs0rY?si=cg-Xy6tdTWqb1xEq" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>`,
								},
							},
						},
					},
				},
				{
					Type:  "section",
					Attrs: map[string]any{"id": "dates"},
					Children: []*Node{
						{
							Type:    "h2",
							Content: "Event Dates",
						},
						{
							Type:    "p",
							Content: "The festival takes place from August 15th to August 18th, 2024.",
						},
					},
				},
				{
					Type:  "section",
					Attrs: map[string]any{"id": "tickets"},
					Children: []*Node{
						{
							Type:    "h2",
							Content: "Tickets",
						},
						{
							Type:    "p",
							Content: "Get your tickets now! Early bird discounts available until June 30th.",
						},
					},
				},
				{
					Type:  "section",
					Attrs: map[string]any{"id": "locations"},
					Children: []*Node{
						{
							Type:    "h2",
							Content: "Locations",
						},
						{
							Type:    "p",
							Content: "The festival will be held at various venues across the city, including open-air stages and indoor theaters.",
						},
					},
				},
				{
					Type:  "section",
					Attrs: map[string]any{"id": "about"},
					Children: []*Node{
						{
							Type:    "h2",
							Content: "About Us",
						},
						{
							Type:    "p",
							Content: "Random Circles is organized by a passionate team dedicated to promoting dance and music culture worldwide.",
						},
					},
				},
			},
		},
	},
}
