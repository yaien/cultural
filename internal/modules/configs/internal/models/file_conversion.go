package models

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

var ErrUnsupportedContentType = errors.New("unsupported content type")

func ConvertFile(dir, inname, outname, contentType string, quality int) (err error) {
	switch {
	case strings.HasPrefix(contentType, "image"):
		err = ConvertImage(dir, inname, outname, contentType, quality)
	case strings.HasPrefix(contentType, "video"):
		err = ConvertVideo(dir, inname, outname, contentType, quality)
	default:
		err = ErrUnsupportedContentType
	}

	return
}

// ConvertImage converts an image file to a different format or quality. it relies on vips for the actual conversion,
// and it is expected to be installed on the system.
// The function takes the directory, input file name, output file name, content type, and quality as parameters.
// It returns an error if the conversion fails.
func ConvertImage(dir, inname, outname, contentType string, width int) (err error) {
	cmd := exec.Command("vips", "thumbnail", path.Join(dir, inname), path.Join(dir, outname)+"[Q=80,strip]", strconv.Itoa(width))
	if err != nil {
		return fmt.Errorf("failed to execute vips command: %w", err)
	}

	var errb bytes.Buffer
	cmd.Stderr = &errb

	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("vips command failed: %s: %w", errb.String(), err)
	}

	return nil
}

func ConvertVideo(dir, inname, outname, contentType string, quality int) (err error) {
	return
}
