package storage

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/lib/worker"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ worker.Handler = (*Handler)(nil)

type Handler struct {
	storage *Storage
}

func NewHandler(s *Storage) *Handler {
	return &Handler{s}
}

func (h *Handler) Handle(ctx context.Context, data map[string]any) error {
	id, ok := data["_id"].(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("invalid file id")
	}

	return h.storage.Convert(ctx, id)
}
