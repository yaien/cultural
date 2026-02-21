package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

var ErrUnsupportedContentType = errors.New("unsupported content type")

func ConvertFile(input io.Reader, output io.Writer, contentType string, variant int) (err error) {
	switch {

	case strings.HasPrefix(contentType, "image"):
		err = ConvertImage(input, output, variant)

	case strings.HasPrefix(contentType, "video"):
		err = ConvertVideo(input, output, variant)

	default:
		err = ErrUnsupportedContentType
	}

	return
}

// ConvertImage converts an image file to a different format or quality. it relies on vips for the actual conversion,
// and it is expected to be installed on the system.
// The function takes the directory, input file name, output file name, content type, and quality as parameters.
// It returns an error if the conversion fails.
func ConvertImage(input io.Reader, output io.Writer, width int) (err error) {
	tempdir, err := os.MkdirTemp("", "vips_")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}

	defer os.RemoveAll(tempdir)

	in, err := os.CreateTemp(tempdir, "input.jpg")
	if err != nil {
		return fmt.Errorf("failed creating input file: %w", err)
	}

	_, err = io.Copy(in, input)
	if err != nil {
		return fmt.Errorf("failed to copy input data to temporary file: %w", err)
	}

	cmd := exec.Command("vips", "thumbnail", in.Name(), path.Join(tempdir, "output.webp")+"[Q=75,strip]", fmt.Sprint(width))
	if err != nil {
		return fmt.Errorf("failed to execute vips command: %w", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("vips command failed: %w, stderr: %s", err, stderr.String())
	}

	out, err := os.OpenInRoot(tempdir, "output.webp")
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}

	_, err = io.Copy(output, out)
	if err != nil {
		return fmt.Errorf("failed to copy output data to output writer: %w", err)
	}

	return nil
}

// ConvertVideo converts a video file to a different format or quality. it relies on ffmpeg for the actual conversion,
func ConvertVideo(input io.Reader, output io.Writer, variant int) (err error) {
	tempdir, err := os.MkdirTemp("", "ffmpeg_")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}

	defer os.RemoveAll(tempdir)

	in, err := os.CreateTemp(tempdir, "input.mp4")
	if err != nil {
		return fmt.Errorf("failed creating input file: %w", err)
	}

	_, err = io.Copy(in, input)
	if err != nil {
		return fmt.Errorf("failed to copy input data to temporary file: %w", err)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", in.Name(),
		"-vf", fmt.Sprintf("scale=-2:%d", variant),
		"-c:v", "libvpx-vp9",
		"-b:v", "0",
		"-crf", "32",
		"-row-mt", "1",
		"-deadline", "good",
		"-c:a", "libopus",
		"-f", "webm",
		path.Join(tempdir, "output.webm"),
	)

	var errp bytes.Buffer
	cmd.Stderr = &errp
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %w, stderr: %s", err, errp.String())
	}

	out, err := os.OpenInRoot(tempdir, "output.webm")
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}

	_, err = io.Copy(output, out)
	if err != nil {
		return fmt.Errorf("failed to copy output data to output writer: %w", err)
	}

	return nil

}
