package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type RenderController struct {
}

func NewRenderController() *RenderController {
	return &RenderController{}
}

type RenderControllerInput struct {
	Map     string                   `json:"map"`
	Key     string                   `json:"key"`
	Layouts map[string]*models.Page  `bson:"layouts" json:"layouts"`
	Fonts   map[string]*models.Font  `bson:"fonts" json:"fonts"`
	Pages   map[string]*models.Page  `bson:"pages" json:"pages"`
	Emails  map[string]*models.Email `bson:"emails" json:"emails"`
	Colors  map[string]string        `bson:"colors" json:"colors"`
}

func (c *RenderController) Render(w http.ResponseWriter, r *http.Request) {

	var input RenderControllerInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteJSONErr(w, models.DecodeError(err))
		return
	}

	switch input.Map {
	case "pages":
		c.RenderPage(w, &input)
	case "layouts":
		c.RenderLayout(w, &input)
	case "emails":
		c.RenderEmail(w, &input)
	default:
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid map: %s", input.Map)))
		return
	}

}

func (c *RenderController) RenderPage(w http.ResponseWriter, input *RenderControllerInput) {
	var buffer bytes.Buffer

	base, err := models.PageTemplate.Clone()
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed decoding template: %w", err))
		return
	}

	page, ok := input.Pages[input.Key]
	if !ok {
		WriteJSONErr(w, &models.Error{Code: "page_not_found", Err: fmt.Errorf("page %q not found in input pages", input.Key)})
		return
	}

	layout, ok := input.Layouts[page.Layout]
	if !ok {
		layout = models.DefaultLayout
	}

	parsed, err := base.Parse(fmt.Sprintf(`{{define "layout_body"}}%s{{end}}{{define "page_body"}}%s{{end}}`, layout.Body, page.Body))
	if err != nil {
		WriteJSONErr(w, &models.Error{Code: "parse_failed", Err: fmt.Errorf("failed parsing template: %w", err)})
		return
	}

	data := models.NewPageData(page, layout).
		WithInlineStyles(true).
		WithInlineScript(true).
		WithColors(input.Colors).
		WithFonts(input.Fonts).
		WithFilePath("/assets/dynamic/files/").
		Data()

	if err := parsed.Execute(&buffer, data); err != nil {
		WriteJSONErr(w, &models.Error{Code: "execution_failed", Err: fmt.Errorf("failed executing template: %w", err)})
		return
	}

	WriteJSON(w, map[string]any{"html": buffer.String()})
}

func (c *RenderController) RenderLayout(w http.ResponseWriter, input *RenderControllerInput) {
	var buffer bytes.Buffer

	base, err := models.PageTemplate.Clone()
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed decoding template: %w", err))
		return
	}

	layout, ok := input.Layouts[input.Key]
	if !ok {
		WriteJSONErr(w, &models.Error{Code: "layout_not_found", Err: fmt.Errorf("layout %q not found in input layouts", input.Key)})
		return
	}

	parsed, err := base.Parse(fmt.Sprintf(`{{define "layout_body"}}%s{{end}}{{define "page_body"}}{{end}}`, layout.Body))
	if err != nil {
		WriteJSONErr(w, &models.Error{Code: "parse_failed", Err: fmt.Errorf("failed parsing template: %w", err)})
		return
	}

	data := models.NewPageData(models.EmptyPage, layout).
		WithInlineStyles(true).
		WithInlineScript(true).
		WithColors(input.Colors).
		WithFonts(input.Fonts).
		WithFilePath("/assets/dynamic/files/").
		Data()

	if err := parsed.Execute(&buffer, data); err != nil {
		WriteJSONErr(w, &models.Error{Code: "execution_failed", Err: fmt.Errorf("failed executing template: %w", err)})
		return
	}

	WriteJSON(w, map[string]any{"html": buffer.String()})

}

func (c *RenderController) RenderEmail(w http.ResponseWriter, input *RenderControllerInput) {
	email, ok := input.Emails[input.Key]
	if !ok {
		WriteJSONErr(w, &models.Error{Code: "email not found", Err: fmt.Errorf("email %s not found in input emails", input.Key)})
		return
	}

	WriteJSON(w, map[string]any{"html": email.Body})
}
