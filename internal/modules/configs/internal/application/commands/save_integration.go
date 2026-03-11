package commands

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SaveIntegrationCommand struct {
	integrations models.IntegrationRepository
}

func NewSaveIntegrationCommand(integrations models.IntegrationRepository) *SaveIntegrationCommand {
	return &SaveIntegrationCommand{
		integrations: integrations,
	}
}

type SaveIntegrationOptions struct {
	OrganizationID primitive.ObjectID
	Name           string
	Data           any
}

func (c *SaveIntegrationCommand) SaveIntegration(ctx context.Context, input SaveIntegrationOptions) error {

	integration, err := c.integrations.Get(ctx, models.GetIntegrationOptions{
		OrganizationID: input.OrganizationID,
		Name:           input.Name,
		Data:           reflect.New(reflect.TypeOf(input).Elem()).Interface(),
	})

	switch {
	case err == nil:
		integration.Data = input.Data
		integration.UpdatedAt = time.Now()
		return c.integrations.Update(ctx, integration)

	case models.IsNotFoundError(err):
		return c.integrations.Create(ctx, &models.Integration{
			ID:             primitive.NewObjectID(),
			OrganizationID: input.OrganizationID,
			Name:           input.Name,
			Data:           input.Data,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})

	default:
		return fmt.Errorf("failed getting integration: %w", err)
	}

}
