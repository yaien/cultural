package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"

	"github.com/yaien/cultural/internal/lib/worker"
)

type Queue = worker.Queue

type Storage struct {
	repository gorm.Interface[File]
	driver     Driver
	queue      *Queue
}

func New(driver Driver, db *gorm.DB, queue *Queue) *Storage {
	return &Storage{gorm.G[File](db), driver, queue}
}

type UploadOptions struct {
	Name           string
	Size           int64
	ContentType    string
	Data           io.Reader
	OrganizationID primitive.ID
}

// Upload uploads a file to the storage and creates a corresponding record in the repository.
func (s *Storage) Upload(ctx context.Context, req *UploadOptions) (file File, err error) {
	count, err := s.repository.
		Where("organization_id = ? and name = ?", req.OrganizationID, req.Name).
		Limit(1).
		Count(ctx, "id")

	if err != nil {
		return file, fmt.Errorf("failed to check file existence: %w", err)
	}

	if count > 0 {
		return file, coderror.New("name_already_exits", errors.New("file already exists"))
	}

	fid := primitive.NewUUID()

	if err = s.driver.Put(fid, req.Size, req.Data); err != nil {
		return file, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	dir, src, err := s.driver.Mount(fid)
	if err != nil {
		return file, fmt.Errorf("failed to mount file: %w", err)
	}

	defer func() {
		if err := s.driver.Unmount(dir); err != nil {
			slog.Error("Failed to unmount file", "error", err)
		}
	}()

	width, height, variant, err := GetDimensionByContentType(ctx, src, req.ContentType)
	if err != nil && !errors.Is(err, ErrUnsupportedContentType) {
		return file, fmt.Errorf("failed to get file dimension: %w", err)
	}

	// Extract preset from content type (e.g., "image/jpeg" -> "image")
	preset := strings.Split(req.ContentType, "/")[0]

	// Remove file extension from name (e.g., "photo.jpg" -> "photo")
	name := req.Name
	if idx := strings.LastIndex(req.Name, "."); idx != -1 {
		name = req.Name[:idx-1]
	}

	file = File{
		OrganizationID: req.OrganizationID,
		Name:           name,
		Preset:         preset,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Formats: []Format{{
			ID:          fid,
			Width:       width,
			Height:      height,
			Variant:     variant,
			Size:        req.Size,
			ContentType: req.ContentType,
		}},
	}

	err = s.repository.Create(ctx, &file)
	if err != nil {
		return file, fmt.Errorf("failed to create file record: %w", err)
	}

	if variant > 0 {
		if err = s.queue.Push(ctx, NewTask(file)); err != nil {
			return file, fmt.Errorf("failed to push compress-file job: %w", err)
		}
	}

	return file, nil
}

// Delete removes a file from the storage and deletes the corresponding record from the repository.
func (s *Storage) Delete(ctx context.Context, organizationID primitive.ID, id primitive.ID) error {
	file, err := s.repository.
		Where("organization_id = ? and id = ?", organizationID, id).
		First(ctx)

	if err != nil {
		return fmt.Errorf("failed to get file from repository: %w", err)
	}

	_, err = s.repository.Where("id = ?", file.ID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete file from repository: %w", err)
	}

	for _, format := range file.Formats {
		err = s.driver.Remove(format.ID)
		if err != nil {
			return fmt.Errorf("failed to delete file from storage: %w", err)
		}
	}

	return nil
}

// Rename updates the name of a file in the repository. It does not change the file in the storage.
func (s *Storage) Rename(ctx context.Context, organizationID primitive.ID, oldName, newName string) error {
	file, err := s.repository.Where("organization_id = ? and name = ?", organizationID, oldName).First(ctx)
	if err != nil {
		return fmt.Errorf("failed to get file from repository: %w", err)
	}

	file.Name = newName
	file.UpdatedAt = time.Now()

	if _, err = s.repository.Updates(ctx, file); err != nil {
		return fmt.Errorf("failed to update file record: %w", err)
	}

	return nil
}

// Get retrieves a file record from the repository by organization ID and name. It does not access the storage.
func (s *Storage) GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ID, name string) (File, error) {
	return s.repository.Where("organization_id = ? and name = ?", organizationID, name).First(ctx)
}

type DownloadOptions struct {
	OrganizationID primitive.ID
	Name           string
	ID             *primitive.ID
	Variant        int
}

type Download struct {
	Format
	Name      string
	Data      io.ReadCloser
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetByOrganizationID retrieves all file records for a given organization ID from the repository. It does not access the storage.
func (s *Storage) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) ([]File, error) {
	return s.repository.Where("organization_id = ?", organizationID).Find(ctx)
}

// Download retrieves a file from the storage and returns it as a Download struct. It first looks up the file record in the repository to get the format information, then retrieves the file data from the storage.
func (s *Storage) Download(ctx context.Context, req *DownloadOptions) (*Download, error) {
	var file File
	var err error

	if req.ID != nil {
		file, err = s.repository.Where("organization_id = ? and id = ?", req.OrganizationID, *req.ID).Take(ctx)
	} else {
		file, err = s.repository.Where("organization_id = ? and name = ?", req.OrganizationID, req.Name).Take(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	format, err := file.GetFormat(req.Variant)
	if err != nil {
		return nil, fmt.Errorf("failed to get file format: %w", err)
	}

	data, err := s.driver.Get(format.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to open file from storage: %w", err)
	}

	res := &Download{
		Format:    format,
		Name:      file.Name,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
		Data:      data,
	}

	return res, nil
}

// Convert checks if the file with the given ID has any missing variants based on its preset. If there are missing variants,
// it mounts the biggest format of the file, performs the necessary conversions to create the missing variants,
// and uploads the converted files back to the storage. Finally,
// it updates the file record in the repository with the new formats and removes the original biggest format if it is dropable.
func (s *Storage) Convert(ctx context.Context, id primitive.ID) error {
	file, err := s.repository.Where("id = ?", id).Take(ctx)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
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

	dir, src, err := s.driver.Mount(state.BiggestFormat.ID)
	if err != nil {
		return fmt.Errorf("failed mounting the biggest format: %w", err)
	}

	defer func() {
		if err := s.driver.Unmount(dir); err != nil {
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

	convertions, err := Convert(ctx, file.Preset, src, outdir)
	if err != nil {
		return fmt.Errorf("failed at preset convert: %w", err)
	}

	for _, conversion := range convertions {

		format := Format{
			ID:          primitive.NewUUID(),
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

		if err = s.driver.Put(format.ID, format.Size, reader); err != nil {
			return fmt.Errorf("failed writing conversion file: %s: %w", conversion.Path, err)
		}

		if err = reader.Close(); err != nil {
			return fmt.Errorf("failed closing reader: %w", err)
		}

		file.Formats = append(file.Formats, format)
	}

	if state.BiggestFormatIsDropable {

		file.Formats = slices.Delete(file.Formats, state.BiggestFormatIndex, state.BiggestFormatIndex+1)

		if err := s.driver.Remove(state.BiggestFormat.ID); err != nil {
			return fmt.Errorf("failed removing original format: %w", err)
		}

	}

	file.UpdatedAt = time.Now()

	if _, err := s.repository.Updates(ctx, file); err != nil {
		return fmt.Errorf("failed updating file: %w", err)
	}

	return nil
}
