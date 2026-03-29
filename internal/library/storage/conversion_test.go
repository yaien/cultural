package storage

import (
	"log/slog"
	"net/http"
	"os"
	"slices"
	"testing"
)

func TestConvertFiles(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		preset   string
		variants []int
	}{
		{"convert a big photo", "testdata/big_photo.jpg", "image", ImagePreset.Variants},
		{"convert a big video", "testdata/big_video.mp4", "video", VideoPreset.Variants},
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

			convertions, err := Convert(ctx, test.preset, test.src, outdir)
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

				_, _, variant, err := GetDimensionByContentType(ctx, conversion.Path, conversion.ContentType)
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
				t.Errorf("expected preset convert function presence '%v', got '%v'", test.state.Preset.Name, state.Preset.Name)
			}

			if !slices.Equal(test.state.MissingVariants, state.MissingVariants) {
				t.Errorf("expected missing variants %v, got %v", test.state.MissingVariants, state.MissingVariants)
			}

		})
	}

}
