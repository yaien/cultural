package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/yaien/cultural/internal/library/worker"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID             primitive.ObjectID `bson:"_id"`
	OrganizationID primitive.ObjectID `bson:"organizationId"`
	Name           string             `bson:"name"`
	Preset         string             `bson:"preset"`
	Hidden         bool               `bson:"hidden"`
	Formats        []Format           `bson:"formats"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

type Format struct {
	ID          primitive.ObjectID `bson:"_id"`
	Variant     int                `bson:"variant"`
	Size        int64              `bson:"size"`
	Width       int                `bson:"width"`
	Height      int                `bson:"height"`
	ContentType string             `bson:"contentType"`
}

type FileRepository interface {
	Create(ctx context.Context, file *File) error
	Update(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*File, error)
	GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) (*File, error)
	GetByOrganizationIDAndID(ctx context.Context, organizationID, id primitive.ObjectID) (*File, error)
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) ([]*File, error)
	RenameByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, oldName, newName string) error
	DeleteByOrganizationIDAndName(ctx context.Context, organizationID primitive.ObjectID, name string) error
}

type FileURLFunc func(filename string, variant ...int) string

func FileURL(filename string, variant ...int) string {
	if len(variant) > 0 {
		return fmt.Sprintf("/assets/dynamic/files/%s?variant=%d", filename, variant[0])
	}
	return fmt.Sprintf("/assets/dynamic/files/%s", filename)
}

// NewExternalFileURLFunc creates a FileURLFunc that generates URLs for files served from the server's external assets endpoint.
func NewExternalFileURLFunc(serverURL string, organizationID primitive.ObjectID) FileURLFunc {
	return func(filename string, variant ...int) string {
		if len(variant) > 0 {
			return fmt.Sprintf("%s/assets/external/%s/%s?variant=%d", serverURL, organizationID.Hex(), filename, variant[0])
		}
		return fmt.Sprintf("%s/assets/external/%s/%s", serverURL, organizationID.Hex(), filename)
	}
}

// GetFormat returns the best format for the given variant.
// If the requested width is less than or equal to 0, or if there is only one format available,
// it returns the biggest format. Otherwise, it finds the nearest bigger or equal format based on the requested width.
func (f *File) GetFormat(v int) (format Format, err error) {
	if len(f.Formats) == 0 {
		err = fmt.Errorf("file has no formats")
		return
	}

	// Sort the formats by their variant (width) in ascending order
	slices.SortFunc(f.Formats, func(a, b Format) int { return a.Variant - b.Variant })

	switch {

	// If the requested width is less than or equal to 0, or if there is only one format available, return the biggest format
	case v <= 0 || len(f.Formats) == 1:
		format = f.Formats[len(f.Formats)-1]

	// default case: find the near bigger or equal format based on the requested width
	default:
		for index := range f.Formats {

			// Find the first format that is smaller than or equal to the requested width
			if v <= f.Formats[index].Variant {
				format = f.Formats[index]
				break
			}

			// If we reached the end of the formats and haven't found a suitable one, return the biggest format
			if index == len(f.Formats)-1 {
				format = f.Formats[index]
				break
			}
		}
	}

	return
}

const GenerateFormatsTaskName = "generate-formats"

// NewGenerateFormatsTask creates a new worker task to generate formats for the file.
func (f *File) NewGenerateFormatsTask() worker.Task {
	return worker.Task{
		Name: GenerateFormatsTaskName,
		Data: map[string]any{"_id": f.ID},
	}
}

// GetFileURL returns the URL for the file's format based on the provided FileURLFunc and variant.
type GetFileDimensionFunc func(ctx context.Context, src string) (width, height, variant int, err error)

// GetFileDimensionByContentType determines the dimensions of a file based on its content type. It supports both images and videos.
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
	defer func() {
		if err := reader.Close(); err != nil {
			slog.Warn("failed to close image file", "error", err)
		}
	}()

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
	Convert     ConvertFunc
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

var ConversionPresets = map[string]*ConversionPreset{
	"image": {
		Convert:     ConvertImage,
		Variants:    []int{320, 640, 1280},
		ContentType: "image/webp",
	},
	"video": {
		Convert:     ConvertVideo,
		Variants:    []int{480, 720},
		ContentType: "video/mp4",
	},
}

var ErrFileHasNoFormats = errors.New("file has no formats")
var ErrFileHasNotConversionPreset = errors.New("file has no conversion preset")

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
