package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

var ErrUnsupportedContentType = errors.New("unsupported content type")

type Conversion struct {
	Variant     int
	Path        string
	ContentType string
	Size        int64
	Width       int
	Height      int
}

type ConvertFunc func(ctx context.Context, src, outdir string, variants []int) (conversions []Conversion, err error)

func ConvertImage(ctx context.Context, src, outdir string, variants []int) (conversion []Conversion, err error) {
	for _, variant := range variants {

		outfile := path.Join(outdir, fmt.Sprintf("output_%d.webp", variant))

		cmd := exec.CommandContext(ctx, "vips", "thumbnail", src, outfile+"[Q=75,strip]", fmt.Sprint(variant))

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("vips command failed: %w, stderr: %s", err, stderr.String())
		}

		stat, err := os.Stat(outfile)
		if err != nil {
			return nil, fmt.Errorf("failed to stat output file: %w", err)
		}

		width, height, _, err := GetImageDimension(ctx, outfile)
		if err != nil {
			return nil, fmt.Errorf("failed getting image dimension: %w", err)
		}

		conversion = append(conversion, Conversion{
			Variant:     variant,
			Path:        outfile,
			ContentType: "image/webp",
			Size:        stat.Size(),
			Width:       width,
			Height:      height,
		})
	}

	return
}

// ConvertVideo converts a video file to a different format or quality. it relies on ffmpeg for the actual conversion,
func ConvertVideo(ctx context.Context, src, outdir string, variants []int) (conversions []Conversion, err error) {
	for _, variant := range variants {

		outfile := path.Join(outdir, fmt.Sprintf("output_%d.webm", variant))

		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-i", src,
			"-vf", fmt.Sprintf("scale=-2:%d", variant),
			"-c:v", "libvpx-vp9",
			"-b:v", "0",
			"-crf", "32",
			"-row-mt", "1",
			"-deadline", "good",
			"-c:a", "libopus",
			"-f", "webm",
			outfile,
		)

		var errp bytes.Buffer
		cmd.Stderr = &errp
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("ffmpeg command failed: %w, stderr: %s", err, errp.String())
		}

		stat, err := os.Stat(outfile)
		if err != nil {
			return nil, fmt.Errorf("failed to stat output file: %w", err)
		}

		width, height, _, err := GetVideoDimension(ctx, outfile)
		if err != nil {
			return nil, fmt.Errorf("failed getting image dimension: %w", err)
		}

		conversions = append(conversions, Conversion{
			Variant:     variant,
			Path:        outfile,
			ContentType: "video/webm",
			Size:        stat.Size(),
			Width:       width,
			Height:      height,
		})
	}

	return conversions, nil
}
