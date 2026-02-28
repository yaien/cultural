package models

import (
	"log/slog"
	"net/http"
	"os"
	"slices"
	"testing"
)

func TestFileGetFormat(t *testing.T) {
	tests := []struct {
		name     string
		variants []int
		query    int
		want     int
		err      string
	}{
		{
			name: "If No Formats: return err",
			err:  "file has no formats",
		},
		{
			name:     "If only one format, return it",
			variants: []int{100},
			query:    300,
			want:     100,
		},
		{
			name:     "If query is 0, return the best format",
			variants: []int{200, 300, 100},
			query:    0,
			want:     300,
		},
		{
			name:     "If query is less than zero formats, return the best format",
			variants: []int{200, 300, 100},
			query:    -1,
			want:     300,
		},
		{
			name:     "If query is greater than all formats, return the best format",
			variants: []int{200, 300, 100},
			query:    400,
			want:     300,
		},
		{
			name:     "If query is between formats, return the best format",
			variants: []int{200, 300, 100},
			query:    250,
			want:     300,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			formats := make([]Format, len(test.variants))
			for i, f := range test.variants {
				formats[i] = Format{Variant: f}
			}

			file := &File{Formats: formats}
			format, err := file.GetFormat(test.query)

			if test.err != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", test.err)
					return
				}

				if err.Error() != test.err {
					t.Errorf("Expected error '%s', got '%s'", test.err, err.Error())
					return
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %s", err.Error())
			}

			if format.Variant != test.want {
				t.Errorf("Expected format quality %d, got %d", test.want, format.Variant)
			}

		})
	}

}

func TestGetFileDimension(t *testing.T) {
	var tests = []struct {
		name        string
		file        string
		contentType string
		width       int
		height      int
		variant     int
		dimension   GetFileDimensionFunc
	}{
		{
			name:        "big_photo",
			file:        "testdata/big_photo.jpg",
			contentType: "image/jpeg",
			width:       3303,
			height:      4954,
			variant:     4954,
			dimension:   GetImageDimension,
		},
		{
			name:        "big_video",
			file:        "testdata/big_video.mp4",
			contentType: "video/mp4",
			width:       1920,
			height:      1080,
			variant:     1080,
			dimension:   GetVideoDimension,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			width, height, quality, err := test.dimension(ctx, test.file)
			if err != nil {
				t.Fatalf("dimension returned an error: %v", err)
			}
			if width != test.width {
				t.Errorf("Expected width %d, got %d", test.width, width)
			}
			if height != test.height {
				t.Errorf("Expected height %d, got %d", test.height, height)
			}
			if quality != test.variant {
				t.Errorf("Expected quality %d, got %d", test.variant, quality)
			}
		})
	}
}

func TestConvertFiles(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		variants  []int
		convert   ConvertFunc
		dimension GetFileDimensionFunc
	}{
		{"convert a big photo", "testdata/big_photo.jpg", []int{320, 640, 1280}, ConvertImage, GetImageDimension},
		{"convert a big video", "testdata/big_video.mp4", []int{420, 720}, ConvertVideo, GetVideoDimension},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

			outdir, err := os.MkdirTemp("", "")
			if err != nil {
				return
			}

			defer func() {
				if err := os.RemoveAll(outdir); err != nil {
					slog.Warn("failed to remove temp directory", "error", err)
				}
			}()

			convertions, err := test.convert(ctx, test.src, outdir, test.variants)
			if err != nil {
				t.Fatalf("failed file convertion: %v", err)
			}

			if len(convertions) != len(test.variants) {
				t.Fatalf("expected %d convertions, got %d", len(test.variants), len(convertions))
			}

			for i, conversion := range convertions {
				if conversion.Variant != test.variants[i] {
					t.Fatalf("expected variant %d, got %d", test.variants[i], conversion.Variant)
				}

				data, err := os.ReadFile(conversion.Path)
				if err != nil {
					t.Fatalf("failed to read convertion file: %v", err)
				}

				detectedContentType := http.DetectContentType(data)
				if detectedContentType != conversion.ContentType {
					t.Fatalf("expected content type %s, got %s", conversion.ContentType, detectedContentType)
				}

				_, _, variant, err := test.dimension(ctx, conversion.Path)
				if err != nil {
					t.Fatalf("failed to get file dimension: %v", err)
				}

				if variant != conversion.Variant {
					t.Fatalf("expected quality %d, got %d", conversion.Variant, variant)
				}
			}

		})
	}
}

func TestFileConversionState(t *testing.T) {
	tests := []struct {
		name    string
		formats []Format
		preset  string
		state   ConversionState

		err string
	}{
		{
			name:    "no formats",
			formats: []Format{},
			err:     "file has no formats",
		},
		{
			name:    "missing all formats",
			preset:  "image",
			formats: []Format{{ContentType: "image/jpeg", Variant: 4954}},
			state: ConversionState{
				BiggestFormat:           Format{ContentType: "image/jpeg", Variant: 4954},
				BiggestFormatIndex:      0,
				BiggestFormatIsDropable: true,
				MissingVariants:         []int{320, 640, 1280},
				Preset:                  ConversionPresets["image"],
			},
		},
		{
			name:   "missing some formats",
			preset: "image",
			formats: []Format{
				{ContentType: "image/jpeg", Variant: 4954},
				{ContentType: "image/webp", Variant: 320},
			},
			state: ConversionState{
				BiggestFormat:           Format{ContentType: "image/jpeg", Variant: 4954},
				BiggestFormatIndex:      0,
				BiggestFormatIsDropable: true,
				MissingVariants:         []int{640, 1280},
				Preset:                  ConversionPresets["image"],
			},
		},
		{
			name:   "all formats present",
			preset: "image",
			formats: []Format{
				{ContentType: "image/webp", Variant: 320},
				{ContentType: "image/webp", Variant: 640},
				{ContentType: "image/webp", Variant: 1280},
				{ContentType: "image/jpeg", Variant: 4954},
			},
			state: ConversionState{
				BiggestFormat:           Format{ContentType: "image/jpeg", Variant: 4954},
				BiggestFormatIndex:      3,
				BiggestFormatIsDropable: true,
				MissingVariants:         []int{},
				Preset:                  ConversionPresets["image"],
			},
		},
		{
			name:   "biggest format is smaller than 1080 but greater than 640",
			preset: "image",
			formats: []Format{
				{ContentType: "image/jpeg", Variant: 920},
				{ContentType: "image/webp", Variant: 640},
			},
			state: ConversionState{
				BiggestFormat:           Format{ContentType: "image/jpeg", Variant: 920},
				BiggestFormatIndex:      0,
				BiggestFormatIsDropable: true,
				MissingVariants:         []int{320},
				Preset:                  ConversionPresets["image"],
			},
		},
		{
			name:   "biggest format is equal to 1080, but its extension is not webp",
			preset: "image",
			formats: []Format{
				{ContentType: "image/jpeg", Variant: 1280},
			},
			state: ConversionState{
				BiggestFormat:           Format{ContentType: "image/jpeg", Variant: 1280},
				BiggestFormatIndex:      0,
				BiggestFormatIsDropable: true,
				MissingVariants:         []int{320, 640, 1280},
				Preset:                  ConversionPresets["image"],
			},
		},
		{
			name:   "biggest format is not dropable because its extension is webp",
			preset: "image",
			formats: []Format{
				{ContentType: "image/webp", Variant: 1280},
			},
			state: ConversionState{
				BiggestFormat:           Format{ContentType: "image/webp", Variant: 1280},
				BiggestFormatIndex:      0,
				BiggestFormatIsDropable: false,
				MissingVariants:         []int{320, 640},
				Preset:                  ConversionPresets["image"],
			},
		},
		{
			name:   "biggest format is dropable because variant is bigger than 1280",
			preset: "image",
			formats: []Format{
				{ContentType: "image/webp", Variant: 3000},
			},
			state: ConversionState{
				BiggestFormat:           Format{ContentType: "image/webp", Variant: 3000},
				BiggestFormatIndex:      0,
				BiggestFormatIsDropable: true,
				MissingVariants:         []int{320, 640, 1280},
				Preset:                  ConversionPresets["image"],
			},
		},
		{
			name:   "file has no conversion preset",
			err:    "file has no conversion preset",
			preset: "application",
			formats: []Format{
				{ContentType: "application/xml", Variant: 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file := &File{Formats: test.formats, Preset: test.preset}
			state, err := file.ConversionState()

			if test.err != "" {
				if err == nil || err.Error() != test.err {
					t.Errorf("expected error '%s', got '%v'", test.err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if state.BiggestFormat.ContentType != test.state.BiggestFormat.ContentType {
				t.Errorf("expected biggest format content type '%s', got '%s'", test.state.BiggestFormat.ContentType, state.BiggestFormat.ContentType)
			}

			if state.BiggestFormat.Variant != test.state.BiggestFormat.Variant {
				t.Errorf("expected biggest format variant '%d', got '%d'", test.state.BiggestFormat.Variant, state.BiggestFormat.Variant)
			}

			if state.BiggestFormatIndex != test.state.BiggestFormatIndex {
				t.Errorf("expected biggest format index '%d', got '%d'", test.state.BiggestFormatIndex, state.BiggestFormatIndex)
			}

			if state.BiggestFormatIsDropable != test.state.BiggestFormatIsDropable {
				t.Errorf("expected biggest format is dropable '%v', got '%v'", test.state.BiggestFormatIsDropable, state.BiggestFormatIsDropable)
			}

			if state.Preset != test.state.Preset {
				t.Errorf("expected preset convert function presence '%v', got '%v'", &test.state.Preset.Convert, &state.Preset.Convert)
			}

			if !slices.Equal(test.state.MissingVariants, state.MissingVariants) {
				t.Errorf("expected missing variants %v, got %v", test.state.MissingVariants, state.MissingVariants)
			}

		})
	}

}
