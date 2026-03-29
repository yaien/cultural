package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"slices"
)

var (
	ErrFileHasNoFormats           = errors.New("file has no formats")
	ErrFileHasNotConversionPreset = errors.New("file has no conversion preset")

	ImagePreset = &ConversionPreset{
		Name:        "image",
		Variants:    []int{320, 640, 1280},
		ContentType: "image/webp",
	}

	VideoPreset = &ConversionPreset{
		Name:        "video",
		Variants:    []int{480, 720},
		ContentType: "video/mp4",
	}

	ConversionPresets = map[string]*ConversionPreset{
		"image": ImagePreset,
		"video": VideoPreset,
	}
)

type Conversion struct {
	Variant     int
	Path        string
	ContentType string
	Size        int64
	Width       int
	Height      int
}

func Convert(ctx context.Context, preset, src, outdir string) (conversions []Conversion, err error) {
	switch preset {
	case "image":
		return ConvertImage(ctx, src, outdir)
	case "video":
		return ConvertVideo(ctx, src, outdir)
	default:
		return nil, fmt.Errorf("unsupported conversion preset: %s", preset)
	}
}

func ConvertImage(ctx context.Context, src, outdir string) (conversion []Conversion, err error) {
	for _, variant := range ImagePreset.Variants {

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
func ConvertVideo(ctx context.Context, src, outdir string) (conversions []Conversion, err error) {
	for _, variant := range VideoPreset.Variants {

		outfile := path.Join(outdir, fmt.Sprintf("output_%d.mp4", variant))

		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-i", src,
			"-vf", fmt.Sprintf("scale=-2:%d", variant),
			"-c:v", "libx264",
			"-preset", "medium",
			"-crf", "23",
			"-movflags", "+faststart",
			"-c:a", "aac",
			"-b:a", "128k",

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
			ContentType: "video/mp4",
			Size:        stat.Size(),
			Width:       width,
			Height:      height,
		})
	}

	return conversions, nil
}

type ConversionPreset struct {
	Name        string
	Variants    []int
	ContentType string
}

type ConversionState struct {
	BiggestFormat           Format
	BiggestFormatIndex      int
	BiggestFormatIsDropable bool
	MissingVariants         []int
	DropableVariantIndexes  []int
	Preset                  *ConversionPreset
}

func (f *File) ConversionState() (state ConversionState, err error) {
	if len(f.Formats) == 0 {
		return state, ErrFileHasNoFormats
	}

	preset, ok := ConversionPresets[f.Preset]
	if !ok {
		return state, ErrFileHasNotConversionPreset
	}

	state.Preset = preset

	already := make(map[int]bool)

	for index, format := range f.Formats {
		if state.BiggestFormat.Variant < format.Variant {
			state.BiggestFormat = format
			state.BiggestFormatIndex = index
			state.BiggestFormatIsDropable = format.Variant > slices.Max(preset.Variants) || format.ContentType != preset.ContentType
		}

		if format.ContentType == preset.ContentType {
			already[format.Variant] = true
		}
	}

	for _, variant := range preset.Variants {
		if already[variant] || variant > state.BiggestFormat.Variant {
			continue
		}

		state.MissingVariants = append(state.MissingVariants, variant)

	}
	return state, nil
}
