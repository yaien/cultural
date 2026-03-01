package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type DraftController struct {
	app *application.Application
}

func NewDraftController(app *application.Application) *DraftController {
	return &DraftController{app: app}
}

func (c *DraftController) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)
	draft, err := c.app.GetDraftByConfigID(ctx, config.ID)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed getting draft for config host %s: %w", config.Host, err))
		return
	}

	WriteJSON(w, draft)
}

func (c *DraftController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	var input struct {
		Fonts   map[string]*models.Font   `json:"fonts"`
		Pages   map[string]*models.Page   `json:"pages"`
		Layouts map[string]*models.Layout `json:"layouts"`
		Emails  map[string]*models.Email  `json:"emails"`
		Colors  map[string]string         `json:"colors"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteJSONErr(w, models.DecodeError(err))
		return
	}

	request := &commands.UpdateDraftRequest{
		ConfigID: config.ID,
		Fonts:    input.Fonts,
		Pages:    input.Pages,
		Layouts:  input.Layouts,
		Emails:   input.Emails,
		Colors:   input.Colors,
	}

	if err := c.app.UpdateDraft(ctx, request); err != nil {
		WriteJSONErr(w, fmt.Errorf("failed updating draft for config host %s: %w", config.Host, err))
		return
	}

	WriteJSONSuccess(w)
}

func (c *DraftController) Commit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(models.ConfigContextKey).(*models.Config)

	if err := c.app.CommitDraft(ctx, config); err != nil {
		WriteJSONErr(w, fmt.Errorf("failed committing draft for config host %s: %w", config.Host, err))
		return
	}

	WriteJSONSuccess(w)
}
