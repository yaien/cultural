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

func (c *RenderController) Render(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Page *models.Page `json:"page"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	if input.Page == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "page is required"})
		return
	}

	var buffer bytes.Buffer

	base, err := models.PageTemplate.Clone()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	parsed, err := base.Parse(fmt.Sprintf(`{{define "body"}}%s{{end}}`, input.Page.Body))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	data := models.NewPageData(config, input.Page).
		WithInlineStyles(true).
		WithFilePath("/assets/dynamic/files/").
		Data()

	err = parsed.Execute(&buffer, data)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"html": buffer.String()})

}
