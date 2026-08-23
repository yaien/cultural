package label

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/lib/cache"
	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

type Draft struct {
	ID        primitive.ID `gorm:"primaryKey;autoIncrement"`
	ConfigID  primitive.ID `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Layouts   map[string]*Layout `gorm:"type:jsonb;serializer:json"`
	Fonts     map[string]*Font   `gorm:"type:jsonb;serializer:json"`
	Pages     map[string]*Page   `gorm:"type:jsonb;serializer:json"`
	Emails    map[string]*Email  `gorm:"type:jsonb;serializer:json"`
	Colors    []*Color           `gorm:"type:jsonb;serializer:json"`
}

type Drafts struct {
	drafts  gorm.Interface[Draft]
	configs gorm.Interface[Config]
	fonts   gorm.Interface[Font]
	cache   *cache.Cache[*Config]
}

func NewDrafts(db *gorm.DB, ch *Cache) *Drafts {
	return &Drafts{
		drafts:  gorm.G[Draft](db),
		configs: gorm.G[Config](db),
		fonts:   gorm.G[Font](db),
		cache:   ch,
	}
}

func (c *Drafts) GetByConfigID(ctx context.Context, configID primitive.ID) (Draft, error) {
	return c.drafts.Where("config_id = ?", configID).First(ctx)
}

func (c *Drafts) CreateColor(ctx context.Context, configID primitive.ID) (*Color, error) {
	draft, err := c.drafts.Where("config_id = ?", configID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	color, err := NewColor(draft.Colors)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tag: %w", err)
	}

	draft.Colors = append(draft.Colors, color)
	draft.UpdatedAt = time.Now()
	if _, err = c.drafts.Updates(ctx, draft); err != nil {
		return nil, fmt.Errorf("failed updating draft: %w", err)
	}

	return color, nil
}

type DraftModelType string

const (
	DraftPageModelType   DraftModelType = "page"
	DraftLayoutModelType DraftModelType = "layout"
	DraftEmailModelType  DraftModelType = "email"
)

type DraftSourceType string

const (
	DraftScriptType DraftSourceType = "script"
	DraftStylesType DraftSourceType = "styles"
	DraftBodyType   DraftSourceType = "body"
)

type UpdateDraftColorOptions struct {
	ConfigID primitive.ID
	ID       primitive.UUID
	Tag      string
	Value    string
}

func (c *Drafts) UpdateColor(ctx context.Context, req *UpdateDraftColorOptions) error {
	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	if req.Tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}

	if req.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	var found bool
	for _, color := range draft.Colors {
		if color.ID == req.ID {
			color.Tag = req.Tag
			color.Value = req.Value
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("color with id %s not found", req.ID)
	}

	draft.UpdatedAt = time.Now()

	if _, err = c.drafts.Updates(ctx, draft); err != nil {
		return fmt.Errorf("failed updating draft: %w", err)
	}

	return nil
}

func (c *Drafts) DeleteColor(ctx context.Context, configID primitive.ID, id primitive.UUID) error {
	draft, err := c.drafts.Where("config_id = ?", configID).First(ctx)
	if err != nil {
		return err
	}

	var deleted bool
	for i, color := range draft.Colors {
		if color.ID == id {
			draft.Colors = append(draft.Colors[:i], draft.Colors[i+1:]...)
			deleted = true
			break
		}
	}

	if !deleted {
		return fmt.Errorf("no color found with id %s", id)
	}

	draft.UpdatedAt = time.Now()

	if _, err = c.drafts.Updates(ctx, draft); err != nil {
		return fmt.Errorf("failed updating draft: %w", err)
	}
	return nil
}

type CreateDraftModelOptions struct {
	ConfigID primitive.ID
	Type     DraftModelType
	Title    string
	Name     string
}

type CreateDraftModelResult struct {
	Draft *Draft
	Model any
}

func (c *Drafts) CreateModel(ctx context.Context, req CreateDraftModelOptions) (*CreateDraftModelResult, error) {
	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	res := CreateDraftModelResult{Draft: &draft}

	switch req.Type {
	case DraftPageModelType:
		_, exists := draft.Pages[req.Name]
		if exists {
			return nil, fmt.Errorf("page with key '%s' already exists", req.Name)
		}

		draft.Pages[req.Name] = &Page{
			Name:  req.Name,
			Title: req.Title,
		}

		res.Model = draft.Pages[req.Name]

	case DraftLayoutModelType:
		_, exists := draft.Layouts[req.Name]
		if exists {
			return nil, fmt.Errorf("layout with key '%s' already exists", req.Name)
		}

		draft.Layouts[req.Name] = &Layout{
			Name:  req.Name,
			Title: req.Title,
		}

		res.Model = draft.Layouts[req.Name]

	default:
		return nil, fmt.Errorf("invalid draft model type: %s", req.Type)
	}

	draft.UpdatedAt = time.Now()

	if _, err := c.drafts.Updates(ctx, draft); err != nil {
		return nil, fmt.Errorf("failed to update draft: %w", err)
	}

	return &res, nil
}

type DeleteDraftModelOptions struct {
	ConfigID primitive.ID
	Type     DraftModelType
	Key      string
}

type DeleteDraftModelResult struct {
	Draft            *Draft
	DefaultModelName string
	DefaultModel     any
}

func (c *Drafts) DeleteModel(ctx context.Context, req DeleteDraftModelOptions) (*DeleteDraftModelResult, error) {
	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	res := DeleteDraftModelResult{Draft: &draft}

	switch req.Type {
	case DraftPageModelType:
		_, exists := draft.Pages[req.Key]
		if !exists {
			return nil, fmt.Errorf("page not found: %s", req.Key)
		}

		if req.Key == DefaultPageName {
			return nil, fmt.Errorf("cannot delete default page")
		}

		delete(draft.Pages, req.Key)

		res.DefaultModelName = DefaultPageName
		res.DefaultModel = draft.Pages[DefaultPageName]

	case DraftLayoutModelType:
		_, exists := draft.Layouts[req.Key]
		if !exists {
			return nil, fmt.Errorf("layout not found: %s", req.Key)
		}

		if req.Key == DefaultLayoutName {
			return nil, fmt.Errorf("cannot delete default layout")
		}

		delete(draft.Layouts, req.Key)

		res.DefaultModelName = DefaultLayoutName
		res.DefaultModel = draft.Layouts[DefaultLayoutName]

	default:
		return nil, fmt.Errorf("invalid draft model type: %s", req.Type)
	}

	draft.UpdatedAt = time.Now()

	if _, err := c.drafts.Updates(ctx, draft); err != nil {
		return nil, fmt.Errorf("failed to update draft: %w", err)
	}

	return &res, nil
}

type UpdateDraftBasicOptions struct {
	ConfigID    primitive.ID
	Type        DraftModelType
	Key         string
	Name        string
	Title       string
	Description string
	Layout      string
	OGImage     string
	OGType      string
	Subject     string
	Preset      string
	PresetParam string
}

func (c *Drafts) UpdateBasic(ctx context.Context, req UpdateDraftBasicOptions) error {
	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	switch req.Type {
	case DraftEmailModelType:
		email, ok := draft.Emails[req.Key]
		if !ok {
			return fmt.Errorf("email not found")
		}

		email.Subject = req.Subject
	case DraftLayoutModelType:
		layout, ok := draft.Layouts[req.Key]
		if !ok {
			return fmt.Errorf("layout not found")
		}

		layout.Name = req.Name
		layout.Title = req.Title

	case DraftPageModelType:
		page, ok := draft.Pages[req.Key]
		if !ok {
			return fmt.Errorf("page not found")
		}

		if req.Key == DefaultPageName && req.Name != DefaultPageName {
			return coderror.Newf("invalid_page_name", "el nombre de la pagina inicial no puede cambiarse")
		}

		if req.Key != req.Name {
			if _, exists := draft.Pages[req.Name]; exists {
				return coderror.Newf("page_already_exist", "ya existe otra pagina con este nombre")
			}

			page.Name = req.Name
			delete(draft.Pages, req.Key)
			draft.Pages[req.Name] = page
		}

		page.Title = req.Title
		page.Description = req.Description
		page.Layout = req.Layout
		page.OGImage = req.OGImage
		page.OGType = req.OGType
		page.Preset = req.Preset
		page.PresetParam = req.PresetParam

	default:
		return fmt.Errorf("invalid type")
	}

	draft.UpdatedAt = time.Now()

	if _, err := c.drafts.Updates(ctx, draft); err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	return nil
}

type UpdateDraftFontOptions struct {
	ConfigID primitive.ID
	Family   string
	Tag      string
}

func (c *Drafts) UpdateFont(ctx context.Context, req UpdateDraftFontOptions) error {
	if req.Tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}

	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get draft by config ID: %w", err)
	}

	font, err := c.fonts.Where("family = ?", req.Family).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get font by family: %w", err)
	}

	if draft.Fonts == nil {
		draft.Fonts = make(map[string]*Font)
	}

	draft.Fonts[req.Tag] = &font
	draft.UpdatedAt = time.Now()

	if _, err := c.drafts.Updates(ctx, draft); err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	return nil
}

type UpdateDraftSourceOptions struct {
	ConfigID   primitive.ID
	Source     string
	ModelType  DraftModelType
	SourceType DraftSourceType
	Key        string
}

func (c *Drafts) UpdateSource(ctx context.Context, req *UpdateDraftSourceOptions) error {
	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return fmt.Errorf("draft not found: %w", err)
	}

	switch req.ModelType {
	case DraftPageModelType:
		page, ok := draft.Pages[req.Key]
		if !ok {
			return fmt.Errorf("page with key '%s' not found", req.Key)
		}

		switch req.SourceType {
		case DraftScriptType:
			page.Script = req.Source
		case DraftStylesType:
			page.Styles = req.Source
		case DraftBodyType:
			page.Body = req.Source
		default:
			return fmt.Errorf("invalid source type: %s", req.ModelType)
		}

	case DraftLayoutModelType:
		layout, ok := draft.Layouts[req.Key]
		if !ok {
			return fmt.Errorf("layout with key '%s' not found", req.Key)
		}

		switch req.SourceType {
		case DraftScriptType:
			layout.Script = req.Source
		case DraftStylesType:
			layout.Styles = req.Source
		case DraftBodyType:
			layout.Body = req.Source
		default:
			return fmt.Errorf("invalid source type: %s", req.ModelType)
		}

	case DraftEmailModelType:
		email, ok := draft.Emails[req.Key]
		if !ok {
			return fmt.Errorf("email with key '%s' not found", req.Key)
		}

		switch req.SourceType {
		case DraftBodyType:
			email.Body = req.Source
		default:
			return fmt.Errorf("invalid source type: %s", req.SourceType)
		}

	default:
		return fmt.Errorf("invalid model type: %s", req.ModelType)
	}

	if _, err := c.drafts.Updates(ctx, draft); err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	return nil
}

type CreateActionRequest struct {
	ConfigID  primitive.ID
	ModelKey  string
	ModelType DraftModelType
	Function  string
}

func (c *Drafts) CreateAction(ctx context.Context, req *CreateActionRequest) error {
	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}
	switch req.ModelType {
	case DraftPageModelType:

		model, ok := draft.Pages[req.ModelKey]
		if !ok {
			return fmt.Errorf("page with name '%s' not found", req.ModelKey)
		}

		var exists bool
		for _, action := range model.Actions {
			if action.Function == req.Function {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("action with function '%s' already exists", req.Function)
		}

		model.Actions = append(model.Actions, &Action{
			Function: req.Function,
		})

		if _, err := c.drafts.Updates(ctx, draft); err != nil {
			return fmt.Errorf("failed to save draft: %w", err)
		}

	case DraftLayoutModelType:

		model, ok := draft.Layouts[req.ModelKey]
		if !ok {
			return fmt.Errorf("page with name '%s' not found", req.ModelKey)
		}

		var exists bool
		for _, action := range model.Actions {
			if action.Function == req.Function {
				exists = true
				break
			}
		}

		if exists {
			return fmt.Errorf("action with function '%s' already exists", req.Function)
		}

		model.Actions = append(model.Actions, &Action{
			Function: req.Function,
		})

		if _, err := c.drafts.Updates(ctx, draft); err != nil {
			return fmt.Errorf("failed to save draft: %w", err)
		}
	}

	return nil
}

type UpdateActionRequest struct {
	ConfigID       primitive.ID
	ModelType      DraftModelType
	ModelKey       string
	Function       string
	TargetFunction string
	Headers        map[string]string
	Body           string
}

func (c *Drafts) UpdateAction(ctx context.Context, req *UpdateActionRequest) error {
	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	switch req.ModelType {
	case DraftPageModelType:

		page, ok := draft.Pages[req.ModelKey]
		if !ok {
			return fmt.Errorf("page with name '%s' not found", req.ModelKey)
		}

		var updated bool
		for _, action := range page.Actions {
			if action.Function == req.TargetFunction {
				action.Body = req.Body
				action.Headers = req.Headers
				action.Function = req.Function
				updated = true
			}
		}

		if !updated {
			return coderror.Newf("function_not_found", "action with function %q not found", req.Function)
		}

		if _, err := c.drafts.Updates(ctx, draft); err != nil {
			return fmt.Errorf("failed to update draft: %w", err)
		}
	case DraftLayoutModelType:
		layout, ok := draft.Layouts[req.ModelKey]
		if !ok {
			return fmt.Errorf("page with name '%s' not found", req.ModelKey)
		}

		var updated bool
		for _, action := range layout.Actions {
			if action.Function == req.TargetFunction {
				action.Body = req.Body
				action.Headers = req.Headers
				action.Function = req.Function
				updated = true
			}
		}

		if !updated {
			return coderror.Newf("function_not_found", "action with function %q not found", req.Function)
		}

		if _, err := c.drafts.Updates(ctx, draft); err != nil {
			return fmt.Errorf("failed to update draft: %w", err)
		}
	}

	return nil
}

type DeleteActionRequest struct {
	ConfigID  primitive.ID
	ModelKey  string
	ModelType DraftModelType
	Function  string
}

func (c *Drafts) DeleteAction(ctx context.Context, req *DeleteActionRequest) error {
	draft, err := c.drafts.Where("config_id = ?", req.ConfigID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	switch req.ModelType {
	case DraftPageModelType:
		page, ok := draft.Pages[req.ModelKey]
		if !ok {
			return fmt.Errorf("page with name '%s' not found", req.ModelKey)
		}

		var deleted bool
		for i, action := range page.Actions {
			if action.Function == req.Function {
				page.Actions = append(page.Actions[:i], page.Actions[i+1:]...)
				deleted = true
				break
			}
		}

		if !deleted {
			return coderror.Newf("action_not_found", "action with function '%s' not found", req.Function)
		}

		if _, err := c.drafts.Updates(ctx, draft); err != nil {
			return fmt.Errorf("failed to update draft: %w", err)
		}
	case DraftLayoutModelType:
		layout, ok := draft.Layouts[req.ModelKey]
		if !ok {
			return fmt.Errorf("layout with key '%s' not found", req.ModelKey)
		}

		var deleted bool
		for i, action := range layout.Actions {
			if action.Function == req.Function {
				layout.Actions = append(layout.Actions[:i], layout.Actions[i+1:]...)
				deleted = true
				break
			}
		}

		if !deleted {
			return coderror.Newf("action_not_found", "action with function '%s' not found", req.Function)
		}

		if _, err := c.drafts.Updates(ctx, draft); err != nil {
			return fmt.Errorf("failed to update draft: %w", err)
		}

	}

	return nil
}

func (c *Drafts) Commit(ctx context.Context, config Config) error {
	draft, err := c.drafts.Where("config_id = ?", config.ID).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	config.Pages = draft.Pages
	config.Emails = draft.Emails
	config.Colors = draft.Colors
	config.Fonts = draft.Fonts
	config.Layouts = draft.Layouts
	config.UpdatedAt = time.Now()

	if _, err := c.configs.Updates(ctx, config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	c.cache.Delete(config.Host)
	return nil
}
