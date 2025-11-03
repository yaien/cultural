package repositories

import (
	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func translate(err error) error {
	switch {
	case err == nil:
		return nil
	case err == mongo.ErrNoDocuments:
		return models.NotFoundError(err)
	default:
		return err
	}
}
