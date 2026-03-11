package integration

import "context"

type InstagramData struct {
}

var Instagram = &Definition{
	Name:        "Instagram",
	Description: "Trae tus posts de instagram a tu web",
	Image:       "instagram.png",
	Path:        "instagram",
	Data:        InstagramData{},
	PageSection: InstagramPage,
	HandleOauth: func(ctx context.Context, oauth *OAuth) (err error) {
		return nil
	},
}
