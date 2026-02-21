package models

import (
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

func GetFileDimension(dir, name, contentType string) (width, height, quality int, err error) {
	switch {

	case strings.HasPrefix(contentType, "image"):
		return GetImageDimension(dir, name)

	case strings.HasPrefix(contentType, "video"):
		return GetVideoDimension(dir, name)

	default:
		err = fmt.Errorf("%w: %s", ErrUnsupportedContentType, contentType)
		return
	}
}

// GetImageDimension decodes the image config to get dimensions without loading the entire image into memory.
func GetImageDimension(dir, name string) (width, height, variant int, err error) {
	var data io.ReadCloser
	data, err = os.OpenInRoot(dir, name)
	if err != nil {
		err = fmt.Errorf("failed getting file from storage: %w", err)
		return
	}

	defer data.Close()

	var img image.Config
	img, _, err = image.DecodeConfig(data)
	if err != nil {
		err = fmt.Errorf("failed decoding image: %w", err)
		return
	}

	width, height, variant = img.Width, img.Height, max(img.Width, img.Height)
	return
}

// GetVideoDimension uses ffprobe to get the dimensions of a video file without loading the entire video into memory.
func GetVideoDimension(dir, name string) (width, height, variant int, err error) {

	filepath := path.Join(dir, name)

	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0",
		filepath)

	var output []byte
	output, err = cmd.Output()
	if err != nil {
		err = fmt.Errorf("failed executing ffprobe: %w", err)
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
