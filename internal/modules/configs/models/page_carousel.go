package models

import (
	"bytes"
	"errors"
	"html/template"
)

var pageCarouselTemplate = template.Must(template.New("carousel").Parse(read("templates/carousel_template.html")))

type pageCarouselItem struct {
	Type string
	Text string
	Src  string
	Alt  string
}

func (p *pageComponents) Carousel(entries ...string) (template.HTML, error) {
	if len(entries)%2 != 0 {
		return "", errors.New("odd number of entries")
	}

	var items []pageCarouselItem

	for i := 0; i < len(entries); i += 2 {
		var item pageCarouselItem
		item.Type = entries[i]
		switch item.Type {
		case "image", "video":
			item.Alt = entries[i+1]
			item.Src = p.options.FilePath + item.Alt
		case "text":
			item.Text = entries[i+1]
		default:
			return "", errors.New("invalid carousel item type")
		}

		items = append(items, item)
	}

	var buff bytes.Buffer

	err := pageCarouselTemplate.Execute(&buff, items)
	if err != nil {
		return "", err
	}

	return template.HTML(buff.String()), nil
}
