package label

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/lib/cache"
	"github.com/yaien/cultural/internal/lib/primitive"
)

type Draft struct {
	ID        primitive.ID `gorm:"primaryKey;autoIncrement"`
	ConfigID  primitive.ID `gorm:"index"`
	Config    *Config
	CreatedAt time.Time
	UpdatedAt time.Time
	Layouts   map[string]*Layout `gorm:"type:jsonb;serializer:json"`
	Fonts     map[string]*Font   `gorm:"type:jsonb;serializer:json"`
	Pages     map[string]*Page   `gorm:"type:jsonb;serializer:json"`
	Emails    map[string]*Email  `gorm:"type:jsonb;serializer:json"`
	Colors    []*Color           `gorm:"type:jsonb;serializer:json"`
}

type DraftRepository interface {
	Update(ctx context.Context, draft *Draft) error
	GetByConfigID(ctx context.Context, id primitive.ID) (*Draft, error)
}

type Drafts struct {
	drafts  DraftRepository
	configs ConfigRepository
	fonts   FontRepository
	cache   *cache.Cache[*Config]
}

func NewDrafts(drafts DraftRepository, configs ConfigRepository, fonts FontRepository, ch *cache.Cache[*Config]) *Drafts {
	return &Drafts{drafts: drafts, configs: configs, fonts: fonts, cache: ch}
}

func (c *Drafts) GetByConfigID(ctx context.Context, configID primitive.ID) (*Draft, error) {
	return c.drafts.GetByConfigID(ctx, configID)
}

func (c *Drafts) CreateColor(ctx context.Context, configID primitive.ID) (*Color, error) {
	draft, err := c.drafts.GetByConfigID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	color, err := NewColor(draft.Colors)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tag: %w", err)
	}

	draft.Colors = append(draft.Colors, color)
	draft.UpdatedAt = time.Now()
	if err = c.drafts.Update(ctx, draft); err != nil {
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
	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
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
		return fmt.Errorf("color with id %d not found", req.ID)
	}

	draft.UpdatedAt = time.Now()

	if err = c.drafts.Update(ctx, draft); err != nil {
		return fmt.Errorf("failed updating draft: %w", err)
	}

	return nil
}

func (c *Drafts) DeleteColor(ctx context.Context, configID primitive.ID, id primitive.UUID) error {
	draft, err := c.drafts.GetByConfigID(ctx, configID)
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
		return fmt.Errorf("no color found with id %d", id)
	}

	draft.UpdatedAt = time.Now()
	return c.drafts.Update(ctx, draft)
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
	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	res := CreateDraftModelResult{Draft: draft}

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
	if err := c.drafts.Update(ctx, draft); err != nil {
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
	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	res := DeleteDraftModelResult{Draft: draft}

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

	if err := c.drafts.Update(ctx, draft); err != nil {
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
}

func (c *Drafts) UpdateBasic(ctx context.Context, req UpdateDraftBasicOptions) error {
	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
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

		if req.Key != DefaultPageName {
			page.Name = req.Name
		}

		page.Title = req.Title
		page.Description = req.Description
		page.Layout = req.Layout
		page.OGImage = req.OGImage
		page.OGType = req.OGType

	default:
		return fmt.Errorf("invalid type")
	}

	draft.UpdatedAt = time.Now()

	if err := c.drafts.Update(ctx, draft); err != nil {
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

	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
	if err != nil {
		return fmt.Errorf("failed to get draft by config ID: %w", err)
	}

	font, err := c.fonts.GetByFamily(ctx, req.Family)
	if err != nil {
		return fmt.Errorf("failed to get font by family: %w", err)
	}

	if draft.Fonts == nil {
		draft.Fonts = make(map[string]*Font)
	}

	draft.Fonts[req.Tag] = font
	draft.UpdatedAt = time.Now()
	if err := c.drafts.Update(ctx, draft); err != nil {
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
	draft, err := c.drafts.GetByConfigID(ctx, req.ConfigID)
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

	if err := c.drafts.Update(ctx, draft); err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	return nil
}

func (c *Drafts) Commit(ctx context.Context, config *Config) error {
	draft, err := c.drafts.GetByConfigID(ctx, config.ID)
	if err != nil {
		return fmt.Errorf("failed to get draft: %w", err)
	}

	config.Pages = draft.Pages
	config.Emails = draft.Emails
	config.Colors = draft.Colors
	config.Fonts = draft.Fonts
	config.Layouts = draft.Layouts
	config.UpdatedAt = time.Now()

	if err := c.configs.Update(ctx, config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	c.cache.Delete(config.Host)
	return nil
}
