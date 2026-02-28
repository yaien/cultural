package controllers

import (
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

	var page, layout *models.Page
	var ok bool

	switch input.Map {
	case "pages":
		page, ok = input.Pages[input.Key]
		if !ok {
			WriteJSONErr(w, &models.Error{Code: "page not found", Err: fmt.Errorf("page %s not found in input pages", input.Key)})
		}
		layout, ok = input.Layouts[page.Layout]
		if !ok {
			layout = models.DefaultLayout
		}
		c.RenderPage(w, page, layout, input.Fonts, input.Colors)

	case "layouts":
		layout, ok = input.Layouts[input.Key]
		if !ok {
			WriteJSONErr(w, &models.Error{Code: "layout not found", Err: fmt.Errorf("layout %s not found in input layouts", input.Key)})
			return
		}

		page = models.EmptyPage
		c.RenderPage(w, page, layout, input.Fonts, input.Colors)
		return

	case "emails":
		c.RenderEmail(w, input.Emails, input.Key)
		return
	default:
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid map: %s", input.Map)))
		return
	}
}

func (c *RenderController) RenderEmail(w http.ResponseWriter, emails map[string]*models.Email, key string) {
	email, ok := emails[key]
	if !ok {
		WriteJSONErr(w, &models.Error{Code: "email not found", Err: fmt.Errorf("email %s not found in input emails", key)})
		return
	}

	WriteJSON(w, map[string]any{"html": email.Body})
}

func (c *RenderController) RenderPage(w http.ResponseWriter, page, layout *models.Page, fonts map[string]*models.Font, colors map[string]string) {
	data := &models.PageData{
		Page:         page,
		Layout:       layout,
		Fonts:        fonts,
		Colors:       colors,
		FilePath:     "/assets/dynamic/files/",
		InlineStyles: true,
		InlineScript: true,
	}

	html, err := models.RenderPage(data)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed to render page: %w", err))
		return
	}
	WriteJSON(w, map[string]any{"html": html})
}
