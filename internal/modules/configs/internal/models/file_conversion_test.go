package models

import (
	"bytes"
	"net/http"
	"os"
	"testing"
)

func TestConvertFile(t *testing.T) {
	tests := []struct {
		name              string
		infile            string
		contentType       string
		outputContentType string
		variant           int
	}{
		{"big photo to 1280", "big_photo.jpg", "image/jpeg", "image/webp", 1280},
		{"big photo to 640", "big_photo.jpg", "image/jpeg", "image/webp", 640},
		{"big photo to 320", "big_photo.jpg", "image/jpeg", "image/webp", 320},
		{"big video to 1080", "big_video.mp4", "video/mp4", "video/webm", 1080},
		{"big video to 720", "big_video.mp4", "video/mp4", "video/webm", 720},
	}

	dir := "testdata"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			input, err := os.OpenInRoot(dir, test.infile)
			if err != nil {
				t.Fatalf("failed to open input file: %v", err)
			}

			defer input.Close()

			var output bytes.Buffer

			err = ConvertFile(input, &output, test.contentType, test.variant)
			if err != nil {
				t.Fatalf("failed file convertion: %v", err)
			}

			detectedContentType := http.DetectContentType(output.Bytes())
			if detectedContentType != test.outputContentType {
				t.Fatalf("expected content type %s, got %s", test.outputContentType, detectedContentType)
			}

			_, _, variant, err := GetFileDimension(&output, test.contentType)
			if err != nil {
				t.Fatalf("failed to get file dimension: %v", err)
			}

			if variant != test.variant {
				t.Fatalf("expected quality %d, got %d", test.variant, variant)
			}

		})
	}
}
