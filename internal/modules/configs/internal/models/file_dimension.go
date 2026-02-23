package models

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

type GetFileDimensionFunc func(ctx context.Context, src string) (width, height, variant int, err error)

func GetFileDimensionByContentType(ctx context.Context, src, contentType string) (width, height, variant int, err error) {
	switch {
	case strings.HasPrefix(contentType, "image/"):
		return GetImageDimension(ctx, src)
	case strings.HasPrefix(contentType, "video/"):
		return GetVideoDimension(ctx, src)
	default:
		err = fmt.Errorf("unsupported content type: %s: %w", contentType, ErrUnsupportedContentType)
		return
	}
}

// GetImageDimension decodes the image config to get dimensions without loading the entire image into memory.
func GetImageDimension(ctx context.Context, src string) (width, height, variant int, err error) {
	reader, err := os.Open(src)
	if err != nil {
		err = fmt.Errorf("failed opening image: %w", err)
		return
	}
	defer reader.Close()

	var img image.Config
	img, _, err = image.DecodeConfig(reader)
	if err != nil {
		err = fmt.Errorf("failed decoding image: %w", err)
		return
	}

	width, height, variant = img.Width, img.Height, max(img.Width, img.Height)
	return
}

// GetVideoDimension uses ffprobe to get the dimensions of a video file without loading the entire video into memory.
func GetVideoDimension(ctx context.Context, src string) (width, height, variant int, err error) {

	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0",
		"-i", src)

	var errp bytes.Buffer
	cmd.Stderr = &errp

	var output []byte
	output, err = cmd.Output()
	if err != nil {
		err = fmt.Errorf("failed executing ffprobe: %s: %w", errp.String(), err)
		return
	}

	dimensions := strings.TrimSpace(string(output))
	parts := strings.Split(dimensions, "x")
	width, err = strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("failed parsing width: %w", err)
		return
	}

	height, err = strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("failed parsing height: %w", err)
		return
	}

	variant = height
	return
}
