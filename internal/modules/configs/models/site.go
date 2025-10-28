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
	Type     string   `bson:"type,omitempty" json:"type,omitempty"`
	Families []string `bson:"families,omitempty" json:"families,omitempty"`
}

type Site struct {
	Styles string `bson:"styles,omitempty" json:"styles,omitempty"`
	Body   *Node  `bson:"body,omitempty" json:"body,omitempty"`
}

var DefaultFonts = &Fonts{
	Type:     "google",
	Families: []string{"Playfair Display", "Inter"},
}

var DefaultSite = &Site{
	Body: &Node{
		Type: "body",
		Style: map[string]string{
			"font-family": "Inter, sans-serif",
		},
		Children: []*Node{
			{
				Type: "header",
				Children: []*Node{
					{
						Type:  "nav",
						Style: map[string]string{"background-color": "#330136"},
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
						Type: "section",
						Children: []*Node{
							{
								Type:  "img",
								Attrs: map[string]any{"src": "/static/landing/brand.png"},
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
								Content: "Random Circles",
								Style:   map[string]string{"font-family": "Playfair Display, sans-serif"},
							},
							{
								Type: "p",
								Content: `Random Circles is a vibrant meeting point for dancers, musicians, and creatives from across the globe - a space where movement, sound, and community come together. Over four days, the festival weaves together dance, live music, workshops, parties, and our legendary BBQ into one immersive celebration rooted in artistic exchange.
                                          From spontaneous ciphers to powerful performances, the atmosphere is charged with creativity and connection. While the dance competition adds energy and expression, the real spirit lies in the shared experience - not in winning.
										  One of the most anticipated moments is the concert, where musicians light up the stage and everyone gathers to celebrate music in its rawest form. And on Sunday, the festival culminates in our signature artistic crossover: an improvisation-driven collaboration between dancers and musicians, unfolding into a unique, shared theater piece.
                                          With over 2,000 guests from nearly 43 countries joining us in past editions, Random Circles continues to grow as a global celebration of culture, movement, and sound.
										  Join us in this extraordinary journey of dance, music, and community spirit!,
										  Follow us on Instagram for the latest updates, videos, and photos.`,
							},
							{
								Type: "carousel",
								Children: []*Node{
									{Type: "div", Content: "random pic 1"},
									{Type: "div", Content: "random pic 2"},
									{Type: "div", Content: "random pic 3"},
								},
							},
							{
								Type:  "section",
								Style: map[string]string{"margin-top": "3px"},
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
	},
}
