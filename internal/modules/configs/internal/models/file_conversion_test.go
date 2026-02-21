package models

import (
	"os"
	"testing"
)

func TestConvertFile(t *testing.T) {
	tests := []struct {
		name        string
		infile      string
		outfile     string
		contentType string
		variant     int
	}{
		{"big photo to 1280", "big_photo.jpg", "big_photo_1280.webp", "image/jpeg", 1280},
		{"big photo to 640", "big_photo.jpg", "big_photo_640.webp", "image/jpeg", 640},
		{"big photo to 320", "big_photo.jpg", "big_photo_320.webp", "image/jpeg", 320},
	}

	dir := "testdata"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := ConvertFile(dir, test.infile, test.outfile, test.contentType, test.variant)
			if err != nil {
				t.Fatalf("failed file convertion: %v", err)
			}

			outfile, err := os.OpenInRoot(dir, test.outfile)
			if err != nil {
				t.Fatalf("failed to open output file: %v", err)
			}

			defer outfile.Close()

			info, err := outfile.Stat()
			if err != nil {
				t.Fatalf("failed to stat output file: %v", err)
			}

			if info.Size() == 0 {
				t.Fatalf("output file is empty")
			}

			_, _, variant, err := GetFileDimension(dir, test.outfile, test.contentType)
			if err != nil {
				t.Fatalf("failed to get file dimension: %v", err)
			}

			if variant != test.variant {
				t.Fatalf("expected quality %d, got %d", test.variant, variant)
			}

		})
	}
}
