package models

import (
	"net/http"
	"os"
	"testing"
)

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

			defer os.RemoveAll(outdir)

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
