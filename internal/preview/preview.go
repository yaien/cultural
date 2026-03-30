package preview

import (
	"context"
	"html/template"
	"maps"

	"github.com/yaien/cultural/internal/coderror"
	"github.com/yaien/cultural/internal/integration"
	"github.com/yaien/cultural/internal/label"
	"github.com/yaien/cultural/internal/storage"
)

type Preview struct {
	registry *integration.Registry
}

func New(registry *integration.Registry) *Preview {
	return &Preview{registry}
}

type GetHTMLOptions struct {
	Key    string
	Type   string
	Config *label.Config
	Draft  *label.Draft
}

func (q *Preview) GetHTML(ctx context.Context, req *GetHTMLOptions) (html string, err error) {
	draft := req.Draft

	switch req.Type {
	case "page":

		page, ok := draft.Pages[req.Key]
		if !ok {
			return "", coderror.Newf(coderror.NotFound, "page %s not found in draft pages", req.Key)
		}

		layout, ok := draft.Layouts[page.Layout]
		if !ok {
			layout = label.DefaultLayout
		}

		return q.renderPage(ctx, &renderPageOptions{
			page:   page,
			layout: layout,
			fonts:  draft.Fonts,
			colors: draft.Colors,
			config: req.Config,
		})

	case "layout":
		layout, ok := draft.Layouts[req.Key]
		if !ok {
			return "", coderror.Newf(coderror.NotFound, "layout %s not found in draft layouts", req.Key)
		}

		page := label.EmptyPage

		return q.renderPage(ctx, &renderPageOptions{
			page:   page,
			layout: layout,
			fonts:  draft.Fonts,
			colors: draft.Colors,
			config: req.Config,
		})

	case "email":
		email, ok := draft.Emails[req.Key]
		if !ok {
			return "", coderror.Newf(coderror.NotFound, "email %s not found in draft emails", req.Key)
		}
		return email.Body, nil
	default:
		return "", coderror.Newf("invalid_type", "invalid type: %s", req.Type)
	}
}

type renderPageOptions struct {
	page   *label.Page
	layout *label.Layout
	fonts  map[string]*label.Font
	colors []*label.Color
	config *label.Config
}

func (q *Preview) renderPage(ctx context.Context, opts *renderPageOptions) (html string, err error) {
	funcs := template.FuncMap{}
	for _, itg := range q.registry.All() {
		if m, ok := itg.(integration.TemplateFuncMap); ok {
			fm := m.TemplateFuncMap(ctx, opts.config)
			maps.Copy(funcs, fm)
		}
	}

	data := &label.PageData{
		Page:                opts.page,
		Layout:              opts.layout,
		Fonts:               opts.fonts,
		Colors:              opts.colors,
		FileURLFunc:         storage.FileURL,
		ExternalFileURLFunc: storage.NewExternalURLFunc(opts.config.Url, opts.config.OrganizationID),
		InlineStyles:        true,
		InlineScript:        true,
		Funcs:               funcs,
	}

	return label.RenderPage(data)
}
