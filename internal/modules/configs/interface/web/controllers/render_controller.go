package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/yaien/cultural/internal/library/render"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type RenderController struct {
}

func NewRenderController() *RenderController {
	return &RenderController{}
}

func (c *RenderController) Render(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Type  string       `json:"type"`
		Page  models.Page  `json:"page"`
		Email models.Email `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var buffer bytes.Buffer

	switch input.Type {
	case "page":
		_ = render.Page(input.Page, nil, render.WithInlineStyles()).Render(r.Context(), &buffer)
	case "email":
		_ = render.Email(input.Email, nil).Render(r.Context(), &buffer)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"html": buffer.String()})

}
