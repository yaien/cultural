package storage

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
)

var _ Storage = (*Local)(nil)

type Local struct {
	root string
}

func NewLocal(root string) *Local {
	return &Local{
		root: root,
	}
}

func (l *Local) Put(id string, size int64, data io.Reader) error {
	root, err := os.OpenRoot(l.root)
	if err != nil {
		return err
	}
	defer root.Close()

	file, err := root.Create(id)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.CopyN(file, data, size)
	if err != nil {
		return err
	}

	return nil
}

func (l *Local) Remove(id string) error {
	root, err := os.OpenRoot(l.root)
	if err != nil {
		return err
	}
	defer root.Close()

	return root.Remove(id)
}

func (l *Local) Get(id string) (io.ReadCloser, error) {
	return os.OpenInRoot(l.root, id)
}

// Dimension returns the width and height of an image or video file without loading the entire file into memory.
// quality is width for images and height for videos, which can be used to determine the quality of the media.
func (l *Local) Dimension(id string, typ string) (width int, height int, quality int, err error) {
	switch {

	// For images, we can decode the config to get dimensions without loading the entire image
	case strings.HasPrefix(typ, "image"):
		var data io.ReadCloser
		data, err = l.Get(id)
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

		width, height, quality = img.Width, img.Height, img.Width

	// For videos, we can use ffprobe to get dimensions without loading the entire video
	case strings.HasPrefix(typ, "video"):

		filepath := path.Join(l.root, id)
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

		quality = height

	default:
		err = fmt.Errorf("unsupported media type: %s", typ)
	}

	return
}
