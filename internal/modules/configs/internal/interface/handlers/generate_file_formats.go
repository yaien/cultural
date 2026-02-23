package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/library/worker"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ worker.Handler = (*GenerateFileFormatHandler)(nil)

type GenerateFileFormatHandler struct {
	files   models.FileRepository
	storage storage.Storage
}

func NewGenerateFileFormatHandler(files models.FileRepository, storage storage.Storage) *GenerateFileFormatHandler {
	return &GenerateFileFormatHandler{
		files:   files,
		storage: storage,
	}
}

func (h *GenerateFileFormatHandler) Handle(ctx context.Context, data map[string]any) error {
	id, ok := data["_id"].(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("invalid file id")
	}

	file, err := h.files.GetByID(ctx, id)
	switch {
	case models.IsNotFoundError(err):
		return nil
	case err != nil:
		return fmt.Errorf("failed to get file: %w", err)
	}

	state, err := file.ConversionState()
	if err != nil {
		return fmt.Errorf("failed to get conversion state: %w", err)
	}

	if len(state.MissingVariants) == 0 {
		return nil
	}

	dir, src, err := h.storage.Mount(state.BiggestFormat.ID.Hex())
	if err != nil {
		return fmt.Errorf("failed mounting the biggest format: %w", err)
	}

	defer func() {
		if err := h.storage.Unmount(dir); err != nil {
			slog.Error("failed unmounting file", "err", err)
		}
	}()

	outdir, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("failed creating temp dir: %w", err)
	}

	defer func() {
		if err := os.RemoveAll(outdir); err != nil {
			slog.Error("failed cleaning outdir", "err", err)
		}
	}()

	convertions, err := state.Preset.Convert(ctx, src, outdir, state.MissingVariants)
	if err != nil {
		return fmt.Errorf("failed at preset convert: %w", err)
	}

	for _, conversion := range convertions {

		format := models.Format{
			ID:          primitive.NewObjectID(),
			ContentType: conversion.ContentType,
			Size:        conversion.Size,
			Height:      conversion.Height,
			Width:       conversion.Width,
			Variant:     conversion.Variant,
		}

		reader, err := os.Open(conversion.Path)
		if err != nil {
			return fmt.Errorf("failed opening conversion file: %s: %w", conversion.Path, err)
		}

		if err = h.storage.Put(format.ID.Hex(), format.Size, reader); err != nil {
			return fmt.Errorf("failed writing conversion file: %s: %w", conversion.Path, err)
		}

		if err = reader.Close(); err != nil {
			return fmt.Errorf("failed closing reader: %w", err)
		}

		file.Formats = append(file.Formats, format)
	}

	if state.BiggestFormatIsDropable {

		file.Formats = slices.Delete(file.Formats, state.BiggestFormatIndex, state.BiggestFormatIndex+1)

		if err := h.storage.Remove(state.BiggestFormat.ID.Hex()); err != nil {
			return fmt.Errorf("failed removing original format: %w", err)
		}

	}

	file.UpdatedAt = time.Now()

	if err := h.files.Update(ctx, file); err != nil {
		return fmt.Errorf("failed updating file: %w", err)
	}

	return nil

}
