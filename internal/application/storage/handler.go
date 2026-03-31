package storage

import (
	"context"
	"fmt"

	"github.com/yaien/cultural/internal/lib/primitive"
	"github.com/yaien/cultural/internal/lib/worker"
)

var _ worker.Handler = (*Handler)(nil)

type Handler struct {
	storage *Storage
}

func NewHandler(s *Storage) *Handler {
	return &Handler{s}
}

func (h *Handler) Handle(ctx context.Context, data map[string]any) error {
	id, ok := data["id"].(primitive.ID)
	if !ok {
		return fmt.Errorf("invalid file id")
	}

	return h.storage.Convert(ctx, id)
}
