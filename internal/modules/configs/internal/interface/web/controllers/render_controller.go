package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type RenderController struct {
}

func NewRenderController() *RenderController {
	return &RenderController{}
}

func (c *RenderController) Render(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Page *models.Page `json:"page"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteJSONErr(w, models.DecodeError(err))
		return
	}

	if input.Page == nil {
		WriteJSONErr(w, models.DecodeError(errors.New("page is required")))
		return
	}

	var buffer bytes.Buffer

	base, err := models.PageTemplate.Clone()
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed decoding template: %w", err))
		return
	}

	parsed, err := base.Parse(fmt.Sprintf(`{{define "body"}}%s{{end}}`, input.Page.Body))
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed parsing template: %w", err))
		return
	}

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	data := models.NewPageData(config, input.Page).
		WithInlineStyles(true).
		WithFilePath("/assets/dynamic/files/").
		Data()

	if err := parsed.Execute(&buffer, data); err != nil {
		WriteJSONErr(w, fmt.Errorf("failed executing template: %w", err))
		return
	}

	WriteJSON(w, map[string]any{"html": buffer.String()})
}
